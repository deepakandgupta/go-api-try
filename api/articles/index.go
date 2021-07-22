package articles

import (
	"log"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/database/mongodbHandler"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const collectionName string = "articles"

type articleWOID struct {
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}

func GetArticles(c *gin.Context) {
	var allArticles []mongodbHandler.Article
	allArticles, err :=  mongodbHandler.GetAllDataFromCollection(collectionName)
	if(err != nil) {
		log.Fatal(err)
		return
	}
	// TODO: Better to convert json and then send
	c.Writer.Header().Set("Eren", "Jaeger")
	c.JSON(http.StatusAccepted, allArticles)
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	article, err := mongodbHandler.GetData(collectionName, "_id", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"id": id})
		return
	}
	c.JSON(200, article)
}

func DeleteArticleByID(c *gin.Context) {
	id := c.Param("id")
	article, err := mongodbHandler.DeleteData(collectionName, "_id", id)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Article Deleted Successfully",
		"data": article,
	})
}

func AddArticleByID(c *gin.Context) {
	var myBodyParams articleWOID
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article, err := mongodbHandler.AddData(collectionName, myBodyParams)
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Article Added Successfully",
		"data": article,
	})
}

func UpdateArticleByID(c *gin.Context) {
	id := c.Param("id")
	var myBodyParams articleWOID
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateValue := bson.M{
			"$set": bson.M{"title": myBodyParams.Title,
			"content": myBodyParams.Content},
	}

	result, err := mongodbHandler.UpdateData(collectionName, "_id", id, updateValue)
	if(err!=nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "Update success",
		"data": result,
	})
}
