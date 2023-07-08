package blackwater

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// To be more specific; fetch connected realms and save them as json files on disk.
func FetchRealmsAndSaveToDisk(r Region, folder string, api *API) {
	api.SetRegion(r, EnUS)
	res, err := api.ConnectedRealmsIndex()
	if err != nil {
		log.Println(err)
		return
	}

	var connectedRealms ConnectedRealmsIndexJson
	json.Unmarshal(res.Body(), &connectedRealms)

	// Get EU and US connected realm indices and save metadata as JSON files in eu/ and us/ folders

	fasthttp.ReleaseResponse(res)

	// Go through every connected realm index
	for _, obj := range connectedRealms.ConnectedRealms {
		s := strings.Split(obj.Href, "/")
		index := strings.Split(s[6], "?")
		crIndex := index[0]
		connectedRealmIndex, err := strconv.Atoi(crIndex)

		if err != nil {
			log.Println(err)
			continue
		}

		// Fetch data on the connected realm
		realmRes, err := api.ConnectedRealm(connectedRealmIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		defer fasthttp.ReleaseResponse(realmRes)

		var connectedRealms ConnectedRealmJson
		json.Unmarshal(realmRes.Body(), &connectedRealms)

		// Create data on the connected realm and its sub-realms
		// and save it on the disk

		fileName := fmt.Sprintf("%s/%s/%d.json", folder, r, connectedRealmIndex)
		f, err := os.Create(fileName)

		if err != nil {
			log.Println(err)
			continue
		}

		enc := json.NewEncoder(f)
		enc.Encode(connectedRealms)
		time.Sleep(100 * time.Millisecond) // Ad-hoc throttle

	}

}
