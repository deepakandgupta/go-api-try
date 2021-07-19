package main

// Import our dependencies. We'll use the standard HTTP library as well as the gorilla router for this app
import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var articles = map[string]Article{
	"1": {
		Title:   "Hello World in GO",
		Content: `fmt.Println("Hello World")`,
	},
	"2": {
		Title:   "Yo in GO",
		Content: `fmt.Println("Yo! From Go!")`,
	},
	"3": {
		Title:   "GOing somewhere with GO",
		Content: `fmt.Println("Vroom! Vroom! in Go!")`,
	},
}

func main() {
	router := gin.Default()


	router.GET("/articles", getArticles)
	router.GET("/article/:id", getArticle)
	router.POST("/article/:id", postArticle)
	router.PUT("/article/:id", updateArticle)
	router.DELETE("/article/:id", deleteArticle)

	log.Fatal(router.Run(":5000"))
}

func getArticles(c *gin.Context) {
	// payload, _ := json.Marshal(articles)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(payload))
	c.JSON(200,articles)
}

func getArticle(c *gin.Context) {
	id := c.Param("id")
	c.Writer.Header().Set("my-key", "Eren Jaeger")
	article, exist := articles[id]
	if exist {
		c.JSON(200, article)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"id": id})
}

func deleteArticle(c *gin.Context) {
	id := c.Param("id")
	article, exist := articles[id]
	if exist {
		c.JSON(200, article)
		delete(articles, id)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"id": id})
	}
}

func postArticle(c *gin.Context) {
	id := c.Param("id")

	var myBodyParams Article
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	articles[id] = myBodyParams
	c.JSON(200, myBodyParams)
}

func updateArticle(c *gin.Context) {
	id := c.Param("id")
	var myBodyParams Article
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exist := articles[id]
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article does not exist"})
		return
	}

	articles[id] = myBodyParams
	c.JSON(200, myBodyParams)
}
