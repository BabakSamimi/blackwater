package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"samimi/blackwater/blackwater"

	_ "github.com/mattn/go-sqlite3"
)

func FileExists(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}

	return nil
}

const (
	databaseFolder = "data/db"
	databaseFile   = databaseFolder + "/bw.db" // This is where all the data will go
	auctionsFolder = "data/auctions"
	realmFolder    = "data/connectedrealms"
	comFolder      = "data/commodities"
)

func main() {

	f, err := os.OpenFile("blackwater.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
	defer f.Close()

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initSql := initCmd.Bool("sql", false, "Sets up the database if it doesn't exist, using sqllite3.")

	realmsCmd := flag.NewFlagSet("realms", flag.ExitOnError)
	euRealms := realmsCmd.Bool("eu", false, "Fetches EU realm indices.")
	usRealms := realmsCmd.Bool("us", false, "Fetches US realm indices.")

	auctionsCmd := flag.NewFlagSet("auctions", flag.ExitOnError)
	euAuctions := auctionsCmd.Bool("eu", false, "Fetches auctions from EU.")
	usAuctions := auctionsCmd.Bool("us", false, "Fetches auctions from US.")

	comCmd := flag.NewFlagSet("commodities", flag.ExitOnError)
	euComs := comCmd.Bool("eu", false, "Fetches commoditiy auctions from EU.")
	usComs := comCmd.Bool("us", false, "Fetches commoditiy auctions from US.")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Expected: blackwater [init|realms|auctions|com] [flags]")
		os.Exit(1)
	}

	os.Mkdir("data", 0777)

	if os.Args[1] == "init" {
		initCmd.Parse(os.Args[2:])
		log.Println("Creating all neccessary folders.")

		os.MkdirAll(auctionsFolder+"/eu", 0777)
		os.Mkdir(auctionsFolder+"/us", 0777)

		os.MkdirAll(comFolder+"/eu", 0777)
		os.MkdirAll(comFolder+"/us", 0777)

		os.MkdirAll(realmFolder+"/eu", 0777)
		os.MkdirAll(realmFolder+"/us", 0777)

		// TODO: SQL init
		if *initSql {
			err = blackwater.InitDB(databaseFile)

			if err != nil {
				log.Print(err)
			}

		}

		os.Exit(0)
	}
	//fmt.Printf("Usage:\n\tblackwater [realms|auctions|commodities] [flags]\n")

	api, apiCreationError := blackwater.NewAPI(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	if apiCreationError != nil {
		log.Fatal(apiCreationError)
	}

	log.Println("Successfully created an API client.")

	switch os.Args[1] {
	case "realms":
		realmsCmd.Parse(os.Args[2:])
		// Fetch realm data
		log.Println("Fetching realm data")

		os.MkdirAll(realmFolder+"/eu", 0777)
		os.MkdirAll(realmFolder+"/us", 0777)

		if *euRealms {
			log.Println("Fetching connected realm indices for EU.")
			blackwater.FetchRealmsAndSaveToDisk(blackwater.EU, realmFolder, api)
			log.Println("Fetching connected realm indices for EU done.")
		}

		if *usRealms {
			log.Println("Fetching connected realm indices for US.")
			blackwater.FetchRealmsAndSaveToDisk(blackwater.US, realmFolder, api)
			log.Println("Fetching connected realm indices for US done.")
		}
	case "auctions":
		// Fetch auctions for all connected realms in us and eu
		os.MkdirAll(auctionsFolder+"/eu", 0777)
		os.MkdirAll(auctionsFolder+"/us", 0777)

		if *euAuctions {
			log.Println("Fetching auctions for EU realms.")
			blackwater.FetchAuctionsAndSaveToDisk(blackwater.EU, auctionsFolder, realmFolder, api)
			log.Println("Fetching of auctions for US realms done.")
		}

		if *usAuctions {
			log.Println("Fetching auctions for US realms.")
			blackwater.FetchAuctionsAndSaveToDisk(blackwater.US, auctionsFolder, realmFolder, api)
			log.Println("Fetching of auctions for US realms done.")
		}

	case "com":
		// Fetch commodities for eu and us
		os.MkdirAll(comFolder+"/eu", 0777)
		os.MkdirAll(comFolder+"/us", 0777)

		if *euComs {
			log.Println("Fetching commodities for EU.")
			blackwater.FetchCommoditiesAndSaveToDisk(blackwater.EU, comFolder, api)
			log.Println("Fetching commodities for EU done.")
		}

		if *usComs {
			log.Println("Fetching commodities for US.")
			blackwater.FetchCommoditiesAndSaveToDisk(blackwater.US, comFolder, api)
			log.Println("Fetching commodities for US done.")
		}

	default:
		fmt.Printf("Expected: blackwater [realms|auctions|commodities] [flags]\n")
		os.Exit(1)
	}

}
