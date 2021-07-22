package mongodbHandler

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName string = "goDemo"

type Article struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}

var MongoClient *mongo.Client

// Functions

func getCollectionAndContext(collectionName string) (context.Context, *mongo.Collection, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := MongoClient.Database(databaseName).Collection(collectionName)
	return ctx, collection, cancel
}

func ConnectToMongoDB() context.Context {
	mongodbPassword := os.Getenv("MONGODB_PASSWORD")
	mongodbDBName := os.Getenv("MONGODB_DB_NAME")
	clientOptions := options.Client().
		ApplyURI(fmt.Sprintf(`mongodb+srv://rozer:%s@cluster0.z1tox.mongodb.net/%s?retryWrites=true&w=majority`, mongodbPassword, mongodbDBName))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return ctx
	}
	MongoClient = client
	fmt.Println("Successfully connected to database")
	return ctx
}

func GetAllDataFromCollection(collectionName string) ([]Article, error) {
	var ds []Article

	ctx, collection, cancel := getCollectionAndContext(collectionName)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return ds, err
	}
	if err = cursor.All(ctx, &ds); err != nil {
		return ds, err
	}
	fmt.Println(ds)
	return ds, nil
}

func GetData(collectionName string, keyToSearch string, value string) (Article, error) {
	var d Article

	ctx, collection, cancel := getCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.SingleResult

	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return d, idErr
		}
		result = collection.FindOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: id},
		})
	} else {
		result = collection.FindOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: value},
		})
	}
	err := result.Decode(&d)
	if err != nil {
		fmt.Println(err)
		return d, err
	}
	fmt.Println(d)
	return d, nil
}

func DeleteData(collectionName string, keyToSearch string, value string) (Article, error) {
	var d Article

	ctx, collection, cancel := getCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.SingleResult

	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return d, idErr
		}
		result = collection.FindOneAndDelete(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: id},
		})
	} else {
		result = collection.FindOneAndDelete(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: value},
		})
	}
	err := result.Decode(&d)
	if err != nil {
		fmt.Println(err)
		return d, err
	}
	fmt.Println(d)
	return d, nil
}

func UpdateData(collectionName string, keyToSearch string, value string, data interface{}) (interface{}, error) {

	ctx, collection, cancel := getCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.UpdateResult
	var err error

	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return result, idErr
		}
		result, err = collection.UpdateOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: id},
		},data)
	} else {
		result, err = collection.UpdateOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: value},
		},data)
	}

	if err != nil {
		fmt.Println(err)
		return result, err
	}
	fmt.Println(result.UpsertedID)
	return result, nil
}

func AddData(collectionName string, d interface{}) (interface{}, error) {

	ctx, collection, cancel := getCollectionAndContext(collectionName)
	defer cancel()

	result, err := collection.InsertOne(ctx, d)

	if err != nil {
		fmt.Println(err)
		return d, err
	}
	fmt.Println(result.InsertedID)
	return d, nil
}
