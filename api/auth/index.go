package auth

import (
	"fmt"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/authController"
	"github.com/deepakandgupta/jwt-auth-noDB/models/authModel"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context){
	var creds authModel.RegistrationCredentials
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := authController.Register(creds)
	
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(status, gin.H{
		"message": "User Registered Successfully",
	})
}

func Login(c *gin.Context) {
	// If the cookie is already present, check if the sessionID is valid
	token := checkCookie(c)
	if token != "" {
		status, _, err := authController.IsAuthenticated(token)
		// If sessionID is valid let user know they are already logged in
		if err == nil && status == http.StatusOK{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Already logged in",
			})
			return
		}
	}

	// Get the parameters from request body
	var creds authModel.Credentials
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Time for the session to expire
	ttlSec := 24*60*60 // 1 day

	// Get token from authController
	status, sessionID, err := authController.Login(creds, ttlSec)

	if err != nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Authorization", sessionID)

	c.JSON(status, bson.M{
		"message": "Logged in succesfully",
	})
}

func Welcome(c *gin.Context) {
	token := checkCookie(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}

	status, name, err := authController.IsAuthenticated(token)
	if err!=nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(status, gin.H{
		"message": fmt.Sprintf("User %s successfully logged in", name),
	})
}

func Logout(c *gin.Context) {
	token := checkCookie(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}
	status, err := authController.Logout(token)

	if err!=nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	// // TODO: change secure to true when using FE instead of postman
	// c.SetCookie(
	// 	"sessionID",
	// 	"",
	// 	-1,
	// 	"/",
	// 	"localhost",
	// 	false,
	// 	true,
	// )

	c.JSON(status, gin.H{
		"message": "User logged out Successfully",
	})
}

func checkCookie(c *gin.Context) string{
	cookie, err := c.Cookie("sessionID")
	if err!=nil{
		return ""
	}
	return cookie
}