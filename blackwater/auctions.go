package blackwater

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

func FetchAuctionsAndSaveToDisk(folder string, realmFolder string, api *API) {

	files, err := ioutil.ReadDir(realmFolder)
	if err != nil {
		log.Printf("%q", err)
	}

	// Every file should be a folder with its name as the connected realm id
	for _, file := range files {

		basename := strings.Split(file.Name(), ".")[0] // no extension

		connectedRealmIndex, err := strconv.Atoi(basename)
		if err != nil {
			log.Printf("%q", err)
			continue
		}

		// Create the AH folder for the realm
		auctionFolder := folder + "/" + basename
		os.Mkdir(auctionFolder, os.ModeAppend)

		res, err := api.Auctions(connectedRealmIndex)
		defer fasthttp.ReleaseResponse(res)

		if err != nil {
			log.Printf("%q", err)
			continue
		}

		lastModified := b2s(res.Header.Peek("Last-Modified"))
		timestamp, err := time.Parse("Mon, _2 Jan 2006 15:04:05 MST", lastModified)
		if err != nil {
			log.Printf("%q", err)
			continue
		}

		timestampUnix := strconv.FormatInt(timestamp.Unix(), 10)
		if err != nil {
			log.Printf("%q", err)
			continue
		}

		// If timestamp doesn't exist, save auction to disk
		fullPath := auctionFolder + "/" + timestampUnix + ".json.gz"
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			log.Printf("%s already exists!", fullPath)
			continue
		}

		f, err := os.Create(fullPath)
		if err != nil {
			log.Printf("%q", err)
			continue
		}

		defer f.Close()
		f.Write(res.Body())

		time.Sleep(100 * time.Millisecond) // Ad-hoc throttle

	}

}

func FetchCommoditiesAndSaveToDisk(folder string, api *API) {

	res, err := api.Commodities()
	defer fasthttp.ReleaseResponse(res)

	if err != nil {
		log.Printf("%q", err)
		return
	}

	lastModified := b2s(res.Header.Peek("Last-Modified"))
	timestamp, err := time.Parse("Mon, _2 Jan 2006 15:04:05 MST", lastModified)
	if err != nil {
		log.Printf("%q", err)
		return
	}

	timestampUnix := strconv.FormatInt(timestamp.Unix(), 10)
	if err != nil {
		log.Printf("%q", err)
		return
	}

	// If timestamp doesn't exist, save auction to disk
	fullPath := folder + "/" + timestampUnix + ".json.gz"
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		log.Printf("%s already exists!", fullPath)
		return
	}

	f, err := os.Create(fullPath)
	if err != nil {
		log.Printf("%q", err)
		return
	}
	defer f.Close()

	f.Write(res.Body())
}
