package main

import (
	"log"

	"github.com/deepakandgupta/jwt-auth-noDB/api/articles"
	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.Default()

	router.GET("/articles", articles.GetArticles)
	router.GET("/article/:id", articles.GetArticle)
	router.POST("/article/:id", articles.PostArticle)
	router.PUT("/article/:id", articles.UpdateArticle)
	router.DELETE("/article/:id", articles.DeleteArticle)

	log.Fatal(router.Run(":5000"))
}
