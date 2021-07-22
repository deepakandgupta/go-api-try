package main

import (
	"log"

	"github.com/deepakandgupta/jwt-auth-noDB/api/articles"
	"github.com/deepakandgupta/jwt-auth-noDB/database/mongodbHandler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const PORT string = ":5000"

func main() {
	router := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connecting to MongoDB
	ctx := mongodbHandler.ConnectToMongoDB()
	defer mongodbHandler.MongoClient.Disconnect(ctx)

	router.GET("/articles", articles.GetArticles)
	router.GET("/article/:id", articles.GetArticleByID)
	router.POST("/article", articles.AddArticleByID)
	router.PUT("/article/:id", articles.UpdateArticleByID)
	router.DELETE("/article/:id", articles.DeleteArticleByID)

	log.Fatal(router.Run(PORT))
}
