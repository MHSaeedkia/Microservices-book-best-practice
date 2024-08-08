package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type News struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title" binding:"required"`
	Content      string             `bson:"content" binding:"required"`
	Author       string             `bson:"author" binding:"required"`
	Created_at   time.Time          `bson:"created_at"`
	Published_at time.Time          `bson:"published_at"`
	News_type    string             `bson:"new_type"`
	Tags         string             `bson:"tags"`
}

func InitializeDB(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func CreateCollection(ctx context.Context, client *mongo.Client, dataBase, collection string) *mongo.Collection {
	return client.Database(dataBase).Collection(collection)
}

func (news *News) Add_news(ctx context.Context, collection *mongo.Collection) error {
	_, err := collection.InsertOne(ctx, news)
	if err != nil {
		return err
	}
	return nil
}

func Add_many_news(ctx context.Context, collection *mongo.Collection, new []interface{}) error {
	_, err := collection.InsertMany(ctx, new)
	if err != nil {
		return err
	}
	return nil
}

func (news *News) Get_single_news(ctx context.Context, collection *mongo.Collection) error {
	filter := bson.D{{Key: "_id", Value: news.ID}}
	err := collection.FindOne(ctx, filter).Decode(&news)
	if err != nil {
		fmt.Println("No news was found with the provided ID")
		return err
	}
	return nil
}

func Get_all_news(ctx context.Context, collection *mongo.Collection, results []News, filter interface{}, number_page, limit string) ([]News, error) {
	numPage, err := strconv.Atoi(number_page)
	if err != nil {
		return nil, err
	}
	if numPage == 0 {
		return nil, fmt.Errorf("number of page should not be 0")
	}
	Limit, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	limitt := int64(Limit)
	skip := (numPage - 1) * int(limitt)

	// Set find options with skip and limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(limitt)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (news *News) Publish_news(ctx context.Context, collection *mongo.Collection) error {
	filter := bson.D{{Key: "_id", Value: news.ID}}
	update := bson.D{
		{"$set", bson.D{{"published_at", time.Now()}}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (news *News) Update_news(ctx context.Context, collection *mongo.Collection) error {
	filter := bson.D{{Key: "_id", Value: news.ID}}
	update := bson.D{
		{"$set", bson.D{{"title", news.Title}}},
		{"$set", bson.D{{"content", news.Content}}},
		{"$set", bson.D{{"author", news.Author}}},
		{"$set", bson.D{{"created_at", news.Created_at}}},
		{"$set", bson.D{{"published_at", news.Published_at}}},
		{"$set", bson.D{{"new_type", news.News_type}}},
		{"$set", bson.D{{"tags", news.Tags}}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (news *News) Delete_news(ctx context.Context, collection *mongo.Collection) error {
	filter := bson.D{{Key: "_id", Value: news.ID}}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func Delete_many_news(ctx context.Context, collection *mongo.Collection) error {
	_, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	return nil
}
