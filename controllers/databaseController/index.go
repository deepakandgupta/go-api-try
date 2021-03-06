package databaseController

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var redisClient *redis.Client


func ConnectToMongoDB() (context.Context, error) {
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")
	mongodbDBName := os.Getenv("MONGODB_DB_NAME")
	mongodbDBUsername := os.Getenv("MONGODB_DB_USERNAME")

	// Connecting to MongoDB
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf(`mongodb+srv://%s:%s@cluster0.z1tox.mongodb.net/%s?retryWrites=true&w=majority`, mongodbDBUsername, mongodbPassword, mongodbDBName))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return ctx, err
	}
	mongoClient = client
	log.Println("Successfully connected to database")
	return ctx, nil
}

func GetMongoClient() *mongo.Client{
	return mongoClient
}

func ConnectToRedisLocalDB () {
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	redisClient = rdb
}

func GetRedisClient() *redis.Client{
	return redisClient
}

func GetCollectionAndContext(collectionName string) (context.Context, *mongo.Collection, context.CancelFunc) {
	mongodbDBName := os.Getenv("MONGODB_DB_NAME")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	collection := mongoClient.Database(mongodbDBName).Collection(collectionName)
	return ctx, collection, cancel
}

