package authController

import (
	"fmt"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/databaseController"
	"github.com/deepakandgupta/jwt-auth-noDB/models/authModel"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const collectionName string = "users"

type Credentials authModel.Credentials

var signedInUsers = make(map[string]string)

func Register(creds Credentials) (int, error){
	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	// check if the user already exist
	result := collection.FindOne(ctx, bson.D{
		bson.E{Key: "username", Value: creds.Username},
	})

	if result != nil {
		err := fmt.Errorf("user already exists")
		return http.StatusForbidden, err
	}

	// hash the password before storing
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if hashingErr!=nil {
		return http.StatusInternalServerError, hashingErr
	}
	
	// store the credentials with hashed password
	var credsToStore = Credentials{
		Username: creds.Username,
		Password: string(hashedPassword),
	}

	_, err := collection.InsertOne(ctx, credsToStore)
	if err != nil {
		return http.StatusInternalServerError, hashingErr
	}

	return http.StatusCreated, nil
}

func Login(creds Credentials) (int, string, error){
	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	var sessionID string

	// get data from database
	result := collection.FindOne(ctx, bson.D{
		bson.E{Key: "username", Value: creds.Username},
	})

	// return if the user does not exists
	if result == nil {
		err := fmt.Errorf("user does not exist")
		return http.StatusNotFound, sessionID, err
	}

	var storedCreds Credentials
	err := result.Decode(&storedCreds)
	if err != nil {
		return http.StatusNotFound, sessionID, err
	}

	// Get the expected password from our database
	expectedPassword:= storedCreds.Password

	err = bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(creds.Password))
	if err!=nil {
		err := fmt.Errorf("wrong password")
		return http.StatusBadRequest, sessionID, err
	}

	// If the password is matched, generate a new session id
	sessionUUID, err := uuid.NewRandom()
	if err!=nil{
		return http.StatusInternalServerError, sessionID, err
	}
	sessionID = sessionUUID.String()

	// store the token in our in memory cache
	signedInUsers[sessionID] = creds.Username
	return	http.StatusCreated ,sessionID, nil
}

func IsAuthenticated(sessionID string) (int, string, error){
	username, exist := signedInUsers[sessionID]

	// is the user session does not exist, return unauthorized
	if !exist{
		err := fmt.Errorf("not authorized")
		return http.StatusUnauthorized, "", err
	}
	return http.StatusOK, username, nil
}

func Logout(sessionID string) (int, error){
	// Logout only if the user is authenticated
	_, _, err:= IsAuthenticated(sessionID)
	if err!=nil{
		return http.StatusBadRequest, err
	}

	// delete user session
	delete(signedInUsers, sessionID)
	return http.StatusOK, nil
}