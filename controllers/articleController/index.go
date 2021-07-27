package articleController

import (
	"fmt"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/databaseController"
	"github.com/deepakandgupta/jwt-auth-noDB/models/articleModel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName string = "articles"


func GetAllArticles() (int, []articleModel.Article, error) {
	var article []articleModel.Article

	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return http.StatusNotFound, article, err
	}
	if err = cursor.All(ctx, &article); err != nil {
		return http.StatusInternalServerError ,article, err
	}
	return http.StatusOK, article, nil
}

func GetArticle(keyToSearch string, value string) (int, articleModel.Article, error) {
	var article articleModel.Article

	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.SingleResult

	// if the search id is key, we need to convert that into mongodb ObjectID before searching
	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return http.StatusBadRequest, article, idErr
		}
		result = collection.FindOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: id},
		})
	} else {
		result = collection.FindOne(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: value},
		})
	}
	if(result.Err() !=nil){
		return http.StatusNotFound, article, result.Err()
	}
	// decode and store in Article model
	err := result.Decode(&article)
	if err != nil {
		fmt.Println(err)
		return http.StatusNotFound, article, err
	}
	return http.StatusOK, article, nil
}

func DeleteArticle(keyToSearch string, value string) (int, articleModel.Article, error) {
	var article articleModel.Article

	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.SingleResult

	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return http.StatusBadRequest, article, idErr
		}
		result = collection.FindOneAndDelete(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: id},
		})
	} else {
		result = collection.FindOneAndDelete(ctx, bson.D{
			bson.E{Key: keyToSearch, Value: value},
		})
	}
	if(result.Err() !=nil){
		return http.StatusNotFound, article, result.Err()
	}

	err := result.Decode(&article)
	if err != nil {
		fmt.Println(err)
		return http.StatusNotFound, article, err
	}
	return http.StatusOK, article, nil
}

func UpdateArticle(keyToSearch string, value string, data interface{}) (int, *mongo.UpdateResult, error) {

	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	var result *mongo.UpdateResult
	var err error

	if keyToSearch == "_id" {
		id, idErr := primitive.ObjectIDFromHex(value)
		if idErr != nil {
			fmt.Println(idErr)
			return http.StatusBadRequest, result, idErr
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
		return http.StatusInternalServerError,result, err
	}
	return http.StatusOK, result, nil
}

func AddArticle(d interface{}) (int, interface{}, error) {
	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	// TODO: check if the id already exist, if it does return article already exists
	result, err := collection.InsertOne(ctx, d)

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError,nil, err
	}
	return http.StatusOK, result.InsertedID, nil
}
