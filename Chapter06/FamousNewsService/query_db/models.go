package query_db

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type QueryNewsModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Version      int                `bson:"version" json:"version" binding:"required"`
	Title        string             `bson:"title" json:"title" binding:"required"`
	Content      string             `bson:"content" json:"content" binding:"required"`
	Author       string             `bson:"author" json:"author" binding:"required"`
	Created_at   time.Time          `bson:"created_at" json:"created_at"`
	Published_at time.Time          `bson:"published_at" json:"published_at"`
	News_type    string             `bson:"new_type" json:"new_type"`
	Tags         string             `bson:"tags" json:"tags"`
}

func InitializeQueryDB(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func (queryNewsModel *QueryNewsModel) RpcQueryConnection() error {
	err := rpc.Register(queryNewsModel)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":4041")
	if err != nil {
		return err
	}
	http.Serve(listener, nil)
	return nil
}

func CreateCollection(ctx context.Context, client *mongo.Client, dataBase, collection string) *mongo.Collection {
	return client.Database(dataBase).Collection(collection)
}

func (queryNewsModel *QueryNewsModel) Get_single_news(ctx context.Context, collection *mongo.Collection) error {
	filter := bson.D{{Key: "_id", Value: queryNewsModel.ID}}
	err := collection.FindOne(ctx, filter).Decode(&queryNewsModel)
	if err != nil {
		return fmt.Errorf("no news was found with the provided ID")
	}
	return nil
}

func Get_all_news(ctx context.Context, collection *mongo.Collection, filter interface{}, number_page, limit string) ([]QueryNewsModel, error) {
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
	queryNewsModel := []QueryNewsModel{}
	err = cursor.All(ctx, &queryNewsModel)
	if err != nil {
		return nil, err
	}
	return queryNewsModel, nil
}
