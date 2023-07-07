package blackwater

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type Client struct {
	ID     string
	Secret string
	Token  *oauth2.Token
}

type Region string
type Locale string
type Namespace string
type Endpoint string

type API struct {
	User       Client
	httpClient *fasthttp.Client
	region     Region
	locale     Locale
}

// Creates a new client
func NewAPI(clientID string, clientSecret string) (api *API, err error) {

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("Client ID or Client Secret was empty")
	}

	api = &API{}
	api.User.ID = clientID
	api.User.Secret = clientSecret

	api.httpClient = &fasthttp.Client{
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	req := fasthttp.AcquireRequest()
	url := fasthttp.AcquireURI()

	// Set URL
	url.Parse(nil, []byte("https://oauth.battle.net/token"))
	url.SetUsername(clientID)
	url.SetPassword(clientSecret)
	req.SetURI(url)
	fasthttp.ReleaseURI(url)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody([]byte("grant_type=client_credentials"))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := fasthttp.AcquireResponse()
	err = api.httpClient.Do(req, res)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Body(), &api.User.Token)
	fasthttp.ReleaseResponse(res)

	if err != nil {
		return nil, err
	}

	// Default
	api.SetRegion(EU, EnUS)

	return
}

const (
	EU Region = "eu"
	US Region = "us"
)

const (
	EnGB Locale = "en_GB"
	EnUS Locale = "en_US"
)

const (
	DynamicEU Namespace = "dynamic-eu"
	DynamicUS Namespace = "dynamic-us"
	StaticEU  Namespace = "static-eu"
	StaticUS  Namespace = "static-us"
)

func (api *API) SetRegion(region Region, locale Locale) {
	api.region = region
	api.locale = locale
}

func (api *API) GetStaticNamespace() Namespace {
	switch api.region {
	case EU:
		return StaticEU
	case US:
		return StaticUS
	default:
		return StaticEU
	}
}

func (api *API) GetDynamicNamespace() Namespace {
	switch api.region {
	case EU:
		return DynamicEU
	case US:
		return DynamicUS
	default:
		return DynamicEU
	}
}

func (api *API) buildUrlDynamic(endpoint string) string {
	return fmt.Sprintf(
		"https://%s.api.blizzard.com/%s?namespace=%s&locale=%s&access_token=%s",
		api.region,
		endpoint,
		api.GetDynamicNamespace(),
		api.locale,
		api.User.Token.AccessToken)
}

func (api *API) buildUrlStatic(endpoint string) string {
	return fmt.Sprintf(
		"https://%s.api.blizzard.com/%s?namespace=%s&locale=%s&access_token=%s",
		api.region,
		endpoint,
		api.GetStaticNamespace(),
		api.locale,
		api.User.Token.AccessToken)
}

// Following types are used for JSON unmarshaling
type ConnectedRealmsIndexJson struct {
	ConnectedRealms []struct {
		Href string `json:"href"`
	} `json:"connected_realms"`
}

type ConnectedRealmJson struct {
	ID     int `json:"id"`
	Realms []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Timezone     string `json:"timezone"`
		IsTournament bool   `json:"is_tournament"`
	} `json:"realms"`
}

type AuctionJson struct {
	Auctions []struct {
		ID   int `json:"id"`
		Item struct {
			ID int `json:"id"`
		} `json:"item"`
		Buyout    int    `json:"buyout"`
		Quantity  int    `json:"quantity"`
		UnitPrice int    `json:"unit_price"`
		TimeLeft  string `json:"time_left"`
	} `json:"auctions"`
}

func (api *API) fetchData(u string) (*fasthttp.Response, error) {

	// Parse the URL
	{
		parsedURL, err := url.Parse(u)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return nil, err
		}

		// Remove the access_token parameter
		q, _ := url.ParseQuery(parsedURL.RawQuery)
		q.Del("access_token")

		// Reconstruct the URL without the access_token parameter
		parsedURL.RawQuery = q.Encode()
		logSafeUrl := parsedURL.String()

		log.Printf("GET on %s\n", logSafeUrl)
	}

	req := fasthttp.AcquireRequest()
	url := fasthttp.AcquireURI()

	url.Parse(nil, []byte(u))
	req.SetURI(url)
	fasthttp.ReleaseURI(url)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Accept", "application/json")

	res := fasthttp.AcquireResponse()
	err := api.httpClient.Do(req, res)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		fasthttp.ReleaseResponse(res)
		return nil, err
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		fasthttp.ReleaseResponse(res)
		return nil, errors.New(res.String())
	}

	return res, nil

}

func (api *API) fetchDataCompressed(u string) (*fasthttp.Response, error) {

	// Parse the URL
	{
		parsedURL, err := url.Parse(u)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return nil, err
		}

		// Remove the access_token parameter
		q, _ := url.ParseQuery(parsedURL.RawQuery)
		q.Del("access_token")

		// Reconstruct the URL without the access_token parameter
		parsedURL.RawQuery = q.Encode()
		logSafeUrl := parsedURL.String()

		log.Printf("GET on %s\n", logSafeUrl)
	}

	req := fasthttp.AcquireRequest()
	url := fasthttp.AcquireURI()

	url.Parse(nil, []byte(u))
	req.SetURI(url)
	fasthttp.ReleaseURI(url)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Accept-Encoding", "gzip")

	res := fasthttp.AcquireResponse()
	err := api.httpClient.Do(req, res)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		fasthttp.ReleaseResponse(res)
		return nil, err
	}

	if res.Header.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(res.String())
	}

	return res, nil

}

func (api *API) ConnectedRealmsIndex() (*fasthttp.Response, error) {

	res, err := api.fetchData(
		api.buildUrlDynamic("data/wow/connected-realm/index"))

	return res, err

}

func (api *API) ConnectedRealm(connectedRealmID int) (*fasthttp.Response, error) {

	res, err := api.fetchData(
		api.buildUrlDynamic(fmt.Sprintf("data/wow/connected-realm/%d", connectedRealmID)))

	return res, err

}

func (api *API) Auctions(connectedRealmID int) (*fasthttp.Response, error) {

	res, err := api.fetchDataCompressed(
		api.buildUrlDynamic(fmt.Sprintf("data/wow/connected-realm/%d/auctions", connectedRealmID)))

	return res, err

}

func (api *API) Commodities() (*fasthttp.Response, error) {

	res, err := api.fetchDataCompressed(
		api.buildUrlDynamic("data/wow/auctions/commodities"))

	return res, err

}
