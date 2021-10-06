package articleModel

import "go.mongodb.org/mongo-driver/bson/primitive"


type Article struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}

type ArticleWOID struct {
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}