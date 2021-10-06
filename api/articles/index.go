package articles

import (
	"fmt"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/articleController"
	"github.com/deepakandgupta/jwt-auth-noDB/controllers/authController"
	"github.com/deepakandgupta/jwt-auth-noDB/models/articleModel"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ArticleWOID articleModel.ArticleWOID

func GetArticles(c *gin.Context) {
	var articles []articleModel.Article
	status, articles, err :=  articleController.GetAllArticles()
	if(err != nil) {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(status, articles)
}

func GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	status, article, err := articleController.GetArticle("_id", id)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	var msg = fmt.Sprintf("Fetched value of id: %s", id)

	var myPayload = gin.H{
		"message": msg,
		"data": article,
	}

	c.JSON(status,myPayload)
}

func DeleteArticleByID(c *gin.Context) {
	if isAuth, _ := checkIfAuthenticated(c) ; !isAuth{
		return
	}

	id := c.Param("id")
	status, article, err := articleController.DeleteArticle("_id", id)
	if(err!=nil){
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	const msg = "Article Deleted Successfully";

	c.JSON(status, gin.H{
		"message": msg,
		"data": article,
	})
}

func AddArticleByID(c *gin.Context) {
	isAuth, name := checkIfAuthenticated(c) ; 
	if (!isAuth){
		return
	}

	var myBodyParams ArticleWOID
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleBody := make(map[string]string)
	articleBody["name"] = name
	articleBody["title"] = myBodyParams.Title
	articleBody["content"] = myBodyParams.Content

	status, article, err := articleController.AddArticle(articleBody)
	if err!=nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	var msg = "Article Added Successfully"

	c.JSON(status, gin.H{
		"message": msg,
		"data": article,
	})
}

func UpdateArticleByID(c *gin.Context) {
	if isAuth, _ := checkIfAuthenticated(c) ; !isAuth{
		return
	}

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
	var msg = fmt.Sprintf("Updated article with id: %s", id)

	c.JSON(status, gin.H{
		"message": msg,
		"data": result,
	})
}

func checkIfAuthenticated(c *gin.Context) (bool, string){
	//  check if the user is authenticated or not
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return false, ""
	}

	status, name, _, err := authController.IsAuthenticated(cookie)
	if err!=nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
			"message": "Login to continue",
		})
		return false, ""
	}
	return true, name
}
