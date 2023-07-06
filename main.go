package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"samimi/blackwater/blackwater"
)

func FileExists(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}

	return nil
}

const (
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

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommands realms, auctions or commodities")
		//os.Exit(1)
	}

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	initSql := initCmd.Bool("sql", false, "Sets up the database if it doesn't exist, using sqllite3.")

	if os.Args[1] == "init" {
		initCmd.Parse(os.Args[2:])
		log.Println("Creating all neccessary folders.")

		os.Mkdir(auctionsFolder, os.ModeAppend)
		os.Mkdir(auctionsFolder+"/eu", os.ModeAppend)
		os.Mkdir(auctionsFolder+"/us", os.ModeAppend)
		os.Mkdir(comFolder+"/eu", os.ModeAppend)
		os.Mkdir(comFolder+"/us", os.ModeAppend)
		os.Mkdir(realmFolder+"/eu", os.ModeAppend)
		os.Mkdir(realmFolder+"/us", os.ModeAppend)

		// TODO: SQL init
		if *initSql {
			log.Println("Setting up SQLLite3 DB.")
		}

		os.Exit(0)
	}

	//realmsCmd := flag.NewFlagSet("realms", flag.ExitOnError)
	//auctionsCmd := flag.NewFlagSet("auctions", flag.ExitOnError)
	//comCmd := flag.NewFlagSet("commodities", flag.ExitOnError)

	api, apiCreationError := blackwater.NewAPI(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	if apiCreationError != nil {
		log.Fatal(apiCreationError)
	}

	fmt.Println("Successfully fetched an Oauth token.")

	os.Mkdir("data", os.ModeAppend)

	switch os.Args[1] {
	case "realms":
		// Fetch realm data
		log.Println("Fetching realm data")

		os.Mkdir(realmFolder, os.ModeAppend)
		os.Mkdir(realmFolder+"/eu", os.ModeAppend)
		os.Mkdir(realmFolder+"/us", os.ModeAppend)

		log.Println("Fetching connected realm indices for EU.")
		api.SetRegion(blackwater.EU, blackwater.EnUS)
		blackwater.FetchRealmsAndSaveToDisk(realmFolder+"/eu", api)
		log.Println("Fetching connected realm indices for EU done.")

		log.Println("Fetching connected realm indices for US.")
		api.SetRegion(blackwater.US, blackwater.EnUS)
		blackwater.FetchRealmsAndSaveToDisk(realmFolder+"/us", api)
		log.Println("Fetching connected realm indices for US done.")

	case "auctions":
		// Fetch auctions for all connected realms in us and eu
		os.Mkdir(auctionsFolder, os.ModeAppend)
		os.Mkdir(auctionsFolder+"/eu", os.ModeAppend)
		os.Mkdir(auctionsFolder+"/us", os.ModeAppend)

		log.Println("Fetching auctions for EU realms.")
		api.SetRegion(blackwater.EU, blackwater.EnUS)
		blackwater.FetchAuctionsAndSaveToDisk(auctionsFolder+"/eu", realmFolder+"/eu", api)
		log.Println("Fetching of auctions for US realms done.")

		log.Println("Fetching auctions for US realms.")
		api.SetRegion(blackwater.US, blackwater.EnUS)
		blackwater.FetchAuctionsAndSaveToDisk(auctionsFolder+"/us", realmFolder+"/us", api)
		log.Println("Fetching of auctions for US realms done.")

	case "com":
		// Fetch commodities for eu and us
		os.Mkdir(comFolder, os.ModeAppend)
		os.Mkdir(comFolder+"/eu", os.ModeAppend)
		os.Mkdir(comFolder+"/us", os.ModeAppend)

		log.Println("Fetching commodities for EU.")
		api.SetRegion(blackwater.EU, blackwater.EnUS)
		blackwater.FetchCommoditiesAndSaveToDisk(comFolder+"/eu", api)
		log.Println("Fetching commodities for EU done.")

		log.Println("Fetching commodities for US.")
		api.SetRegion(blackwater.US, blackwater.EnUS)
		blackwater.FetchCommoditiesAndSaveToDisk(comFolder+"/us", api)
		log.Println("Fetching commodities for US done.")

	}

}
