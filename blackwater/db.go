package blackwater

import (
	"database/sql"
	"log"
)

func InitDB(dbPath string) error {
	log.Println("Setting up SQLLite3 DB.")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Print("Created a handle to the DB.")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS realms(realmID INTEGER NOT NULL PRIMARY KEY,
		connectedRealmID INTEGER,
		region INTEGER,
		name TEXT,
		timezone TEXT);`)

	if err != nil {
		return err
	}

	log.Println("Created Realm table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS processedFiles(fileHash INTEGER NOT NULL PRIMARY KEY,
		bytes INTEGER);`)

	if err != nil {
		return err
	}

	log.Println("Created Processed Files table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items(itemID INTEGER NOT NULL PRIMARY KEY,
		itemClassID INTEGER,
		itemSubclassID INTEGER,
		commodity INTEGER,
		name TEXT);`)

	if err != nil {
		return err
	}

	log.Println("Created Items table")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS auctions(id INTEGER NOT NULL PRIMARY KEY,
		itemID INTEGER,
		buyout INTEGER,
		quantity INTEGER,
		unitPrice INTEGER,
		timeLeft TEXT,
		connectedRealmID INTEGER);`)

	if err != nil {
		return err
	}

	log.Println("Created Auctions table")

	return nil
}
