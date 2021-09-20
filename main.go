package main

import (
	"log"
	"time"

	"github.com/deepakandgupta/jwt-auth-noDB/api/articles"
	"github.com/deepakandgupta/jwt-auth-noDB/api/auth"
	"github.com/deepakandgupta/jwt-auth-noDB/controllers/databaseController"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const PORT string = ":5000"

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET","POST","PUT", "PATCH","OPTIONS","DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connecting to MongoDB
	ctxMongo, err := databaseController.ConnectToMongoDB()
	if err!=nil{
		log.Fatal("Cannot connect to MongoDB Database")
	}
	mongoClient := databaseController.GetMongoClient()
	defer mongoClient.Disconnect(ctxMongo)
	
	// Connecting to Redis
	databaseController.ConnectToRedisLocalDB()

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.GET("/dashboard", auth.Welcome)
	router.GET("/logout", auth.Logout)

	router.GET("/articles", articles.GetArticles)
	router.GET("/article/:id", articles.GetArticleByID)
	router.POST("/article", articles.AddArticleByID)
	router.PUT("/article/:id", articles.UpdateArticleByID)
	router.DELETE("/article/:id", articles.DeleteArticleByID)

	log.Fatal(router.Run(PORT))
}
