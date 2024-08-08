package main

import (
	"context"
	"log"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	ctx        context.Context
	collection *mongo.Collection
	Router     *gin.Engine
}

func (a *App) Initialize(client *mongo.Client, dataBase, collection string) {
	a.ctx = context.TODO()
	a.collection = CreateCollection(a.ctx, client, dataBase, collection)
	gin.SetMode(gin.ReleaseMode)
	a.Router = gin.Default()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.GET("/famous/news/:news_id", a.get_single_news)
	a.Router.GET("/famous/news/all/:num_page/:limit", a.get_all_news)
	a.Router.POST("/famous/news", a.add_news)
	a.Router.GET("/famous/news/:news_id/publish/", a.publish_news)
	a.Router.PUT("/famous/news", a.update_news)
	a.Router.DELETE("/famous/news/:news_id", a.delete_news)
}

func (a *App) Run(addr string) {
	log.Fatal(a.Router.Run(addr))
}

func (a *App) get_single_news(c *gin.Context) {
	news_id := c.Param("news_id")
	// Regular expression to match numeric values
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(news_id) {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}

	news := News{}
	news, err := idSetter(news_id, news)
	if err != nil {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}
	err = news.Get_single_news(a.ctx, a.collection)
	if err != nil {
		respondWithError(c, 400, err.Error())
		return
	}
	respondWithJSON(c, 200, news, "Success")
}

func (a *App) get_all_news(c *gin.Context) {
	num_page := c.Param("num_page")
	limit := c.Param("limit")
	// Regular expression to match numeric values
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(num_page) || !regex.MatchString(limit) {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}
}

func (a *App) add_news(c *gin.Context) {

}

func (a *App) publish_news(c *gin.Context) {
	news_id := c.Param("news_id")
	// Regular expression to match numeric values
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(news_id) {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}
}

func (a *App) update_news(c *gin.Context) {
	news_id := c.Param("news_id")
	// Regular expression to match numeric values
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(news_id) {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}
}

func (a *App) delete_news(c *gin.Context) {
	news_id := c.Param("news_id")
	// Regular expression to match numeric values
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if !regex.MatchString(news_id) {
		respondWithError(c, 400, "StatusBadRequest")
		return
	}
}

func respondWithError(c *gin.Context, code int, message string) {
	respondWithJSON(c, code, message, "Error")
}

func respondWithJSON(c *gin.Context, code int, message interface{}, status string) {
	c.JSON(code, gin.H{
		"Status":  status,
		"Code":    code,
		"Message": message,
	})
}

func idSetter(id string, news News) (News, error) {
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return News{}, err
	}
	news.ID = ID
	return news, nil
}
