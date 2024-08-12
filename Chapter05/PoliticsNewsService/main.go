package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
)

type CONFIG struct {
	debug      bool
	dbSettings MONGODB_SETTINGS
}

type MONGODB_SETTINGS struct {
	db         string
	host       string
	collection string
}

func main() {
	config := CONFIG{}

	debugValue, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&config.debug, "debug", debugValue, "If it true , we are in Development -- If it is false , we are in Production")
	flag.StringVar(&config.dbSettings.db, "db_name", os.Getenv("APP_DATABASE_NAME"), "Database name")
	flag.StringVar(&config.dbSettings.host, "db_host", os.Getenv("APP_DATABASE_HOST"), "Database host")
	flag.StringVar(&config.dbSettings.collection, "db_collection", os.Getenv("APP_DATABASE_COLLECTION"), "Database collection")

	client, err := InitializeDB(context.TODO(), config.dbSettings.host)
	if err != nil {
		panic(err)
	}

	a := App{}
	a.Initialize(client, config.dbSettings.db, config.dbSettings.collection)
	a.Run(":5000")
}
