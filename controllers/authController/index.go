package authController

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/databaseController"
	"github.com/deepakandgupta/jwt-auth-noDB/helpers"
	"github.com/deepakandgupta/jwt-auth-noDB/models/authModel"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const collectionName string = "users"

var ctxRedis = context.Background()


func Register(creds authModel.RegistrationCredentials) (int, error){
	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	// get username to check if it is valid and whether the user has already registred or not
	creds.Username = strings.ToLower(creds.Username)
	creds.Name = strings.ToLower(creds.Name)

	emailErr := helpers.CheckValidEmail(creds.Username)
	if emailErr!=nil {
		return http.StatusBadRequest, emailErr
	}
	
	// Password should be following certain standard
	isValidPass := helpers.CheckIfValidPassword(creds.Password)
	if !isValidPass {
		return http.StatusBadRequest, fmt.Errorf("password not according policy")
	}

	isValidName := helpers.CheckIfValidName(creds.Name)
	if !isValidName {
		return http.StatusBadRequest, fmt.Errorf("name not according policy")
	}

	// check if the user already exist
	result := collection.FindOne(ctx, bson.D{
		bson.E{Key: "username", Value: creds.Username},
	})
	
	if result.Err() == nil {
		err := fmt.Errorf("user already exists")
		return http.StatusForbidden, err
	}

	// hash the password before storing
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if hashingErr!=nil {
		return http.StatusInternalServerError, hashingErr
	}
	
	// store the credentials with hashed password
	var credsToStore = authModel.RegistrationCredentials{
		Name: creds.Name,
		Username: creds.Username,
		Password: string(hashedPassword),
	}

	_, err := collection.InsertOne(ctx, credsToStore)
	if err != nil {
		return http.StatusInternalServerError, hashingErr
	}

	return http.StatusCreated, nil
}

func Login(creds authModel.Credentials, ttlSec int) (int, string, string, string, error){
	ctx, collection, cancel := databaseController.GetCollectionAndContext(collectionName)
	defer cancel()

	var username string
	var name string
	var sessionID string

	// Check if username if valid
	creds.Username = strings.ToLower(creds.Username)

	emailErr := helpers.CheckValidEmail(creds.Username)

	if emailErr!=nil {
		return http.StatusBadRequest, sessionID, name, username, emailErr
	}

	// get data from database
	result := collection.FindOne(ctx, bson.D{
		bson.E{Key: "username", Value: creds.Username},
	})

	// return if the user does not exists
	if result.Err() != nil {
		err := fmt.Errorf("user does not exist")
		return http.StatusNotFound, sessionID, name, username, err
	}

	var storedCreds authModel.RegistrationCredentials
	err := result.Decode(&storedCreds)
	if err != nil {
		return http.StatusNotFound, sessionID, name, username, err
	}

	// Get the expected password from our database
	expectedPassword:= storedCreds.Password

	err = bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(creds.Password))
	if err!=nil {
		err := fmt.Errorf("invalid username or password")
		return http.StatusBadRequest, sessionID, name, username, err
	}

	// If the password is matched, generate a new session id
	sessionUUID, err := uuid.NewRandom()
	if err!=nil{
		return http.StatusInternalServerError, sessionID, name, username, err
	}
	sessionID = sessionUUID.String()
	var rdb = databaseController.GetRedisClient()
	// store the token in our in memory cache
	redisVal := helpers.ToGOB64(storedCreds.Name, storedCreds.Username)
	err = rdb.Set(ctxRedis, sessionID, redisVal, time.Duration(ttlSec)*time.Second).Err()
	if err != nil {
		log.Fatal("Redis Error: Cannot set key")
	}
	
	name = storedCreds.Name
	username = storedCreds.Username
	
	return	http.StatusCreated ,sessionID, name, username, nil
}

func IsAuthenticated(sessionID string) (int, string, string, error){
	var rdb = databaseController.GetRedisClient()
	value, err := rdb.Get(ctxRedis, sessionID).Result()
	data := helpers.FromGOB64(value)
	// is the user session does not exist, return unauthorized
	if err == redis.Nil{
		err := fmt.Errorf("not authorized")
		return http.StatusUnauthorized, "", "", err
	} else if err!=nil{
		return http.StatusInternalServerError, "", "", err
	}
	return http.StatusOK, data.Name, data.Username, nil
}

func Logout(sessionID string) (int, error){
	// Logout only if the user is authenticated
	_, _, _, err:= IsAuthenticated(sessionID)
	if err!=nil{
		return http.StatusBadRequest, err
	}
	var rdb = databaseController.GetRedisClient()
	// delete user session
	_, err = rdb.Del(ctxRedis, sessionID).Result()
	if err!=nil{
		log.Print("Cannot logout user, try again")
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}