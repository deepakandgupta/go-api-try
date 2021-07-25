package articles

import (
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/articleController"
	"github.com/deepakandgupta/jwt-auth-noDB/models/articleModel"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ArticleWOID articleModel.ArticleWOID

func GetArticles(c *gin.Context) {
	var articles []articleController.Article
	status, articles, err :=  articleController.GetAllArticles()
	if(err != nil) {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Eren", "Jaeger")
	c.JSON(status, articles)
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	status, article, err := articleController.GetArticle("_id", id)
	if err != nil {
		c.JSON(status, gin.H{"id": id})
		return
	}
	c.JSON(status, article)
}

func DeleteArticleByID(c *gin.Context) {
	id := c.Param("id")
	status, article, err := articleController.DeleteArticle("_id", id)
	if(err!=nil){
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{
		"message": "Article Deleted Successfully",
		"data": article,
	})
}

func AddArticleByID(c *gin.Context) {
	var myBodyParams ArticleWOID
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status, article, err := articleController.AddArticle(myBodyParams)
	if err!=nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{
		"message": "Article Added Successfully",
		"data": article,
	})
}

func UpdateArticleByID(c *gin.Context) {
	id := c.Param("id")
	var myBodyParams ArticleWOID
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: use replaceOne instead of updateOne if updating the whole article
	updateValue := bson.M{
			"$set": bson.M{"title": myBodyParams.Title,
			"content": myBodyParams.Content},
	}

	status, result, err := articleController.UpdateArticle("_id", id, updateValue)
	if(err!=nil){
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{
		"data": result,
	})
}
