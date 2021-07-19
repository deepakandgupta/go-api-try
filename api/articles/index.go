package articles

import (
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

func GetArticles(c *gin.Context) {
	// payload, _ := json.Marshal(articles)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(payload))
	c.JSON(200,articles)
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	c.Writer.Header().Set("my-key", "Eren Jaeger")
	article, exist := articles[id]
	if exist {
		c.JSON(200, article)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"id": id})
}

func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	article, exist := articles[id]
	if exist {
		c.JSON(200, article)
		delete(articles, id)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"id": id})
	}
}

func PostArticle(c *gin.Context) {
	id := c.Param("id")

	var myBodyParams Article
	if err := c.ShouldBindJSON(&myBodyParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	articles[id] = myBodyParams
	c.JSON(200, myBodyParams)
}

func UpdateArticle(c *gin.Context) {
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
