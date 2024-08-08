package main

import (
	"context"
)

const (
	DBNAME     = "testing" //"FamouseNewsService"
	COLLECTION = "News"
)

func main() {

	connectionString := "mongodb://localhost:27017" //os.Getenv("DATABASE_URL")
	client, err := InitializeDB(context.TODO(), connectionString)
	if err != nil {
		panic(err)
	}

	a := App{}
	a.Initialize(client, DBNAME, COLLECTION)
	a.Run(":5000")
}
