package authModel

// Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}