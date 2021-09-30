package authModel

// Create a struct to read the username and password from the request body
type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Registration parameters can vary a lot from login parameters
// Create a struct to read the name, username and password from the request body
type RegistrationCredentials struct {
	Name string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}