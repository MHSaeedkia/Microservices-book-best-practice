package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/MHSaeedkia/microservice/FamousNewsService/query_db"
)

type CONFIG struct {
	debug             bool
	queryDbSettings   QUERY_DB_SETTINGS
	commandDbSettings COMMAND_DB_SETTINGS
}

type QUERY_DB_SETTINGS struct {
	db         string
	host       string
	collection string
}

type COMMAND_DB_SETTINGS struct {
	uri string
}

func main() {
	config := CONFIG{}

	debugValue, err := strconv.ParseBool(os.Getenv("APP_DEBUG"))
	if err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&config.debug, "debug", debugValue, "If it true , we are in Development -- If it is false , we are in Production")
	flag.StringVar(&config.queryDbSettings.db, "db_query_name", os.Getenv("APP_DATABASE_NAME"), "Database name")
	flag.StringVar(&config.queryDbSettings.host, "db_query_host", os.Getenv("APP_DATABASE_HOST"), "Database host")
	flag.StringVar(&config.queryDbSettings.collection, "db_query_collection", os.Getenv("APP_DATABASE_COLLECTION"), "Database collection")
	flag.StringVar(&config.commandDbSettings.uri, "db_command", os.Getenv("APP_DATABASE_COMMAND"), "Database command")

	client, err := query_db.InitializeQueryDB(context.TODO(), config.queryDbSettings.host)
	if err != nil {
		panic(err)
	}

	// connection, err := command_db.InitializeCommandDB(config.commandDbSettings.uri)
	// if err != nil {
	// 	panic(err)
	// }

	queryNewsModel := query_db.QueryNewsModel{}
	go func() {
		err := queryNewsModel.RpcQueryConnection()
		if err != nil {
			panic(err)
		}
	}()

	a := App{}
	a.Initialize(client, config.queryDbSettings.db, config.queryDbSettings.collection)
	a.Run(":5000")
}
