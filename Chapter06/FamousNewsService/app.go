package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	ctx    context.Context
	Router *gin.Engine
}

func (a *App) Initialize(client *mongo.Client, dataBase, collection string) {
	a.ctx = context.TODO()
	// a.collection = CreateCollection(a.ctx, client, dataBase, collection)
	gin.SetMode(gin.ReleaseMode)
	a.Router = gin.Default()
}

func (a *App) Run(addr string) {
	log.Fatal(a.Router.Run(addr))
}
