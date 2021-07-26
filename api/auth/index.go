package auth

import (
	"fmt"
	"net/http"

	"github.com/deepakandgupta/jwt-auth-noDB/controllers/authController"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context){
	var creds authController.Credentials
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
	cookie, cookieReadErr := c.Cookie("sessionID")
	if cookieReadErr == nil {	
		status, _, err := authController.IsAuthenticated(cookie)
		// If sessionID is valid let user know they are already logged in
		if err == nil && status == http.StatusOK{
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Already logged in",
			})
			return
		}
	}

	// Get the parameters from request body
	var creds authController.Credentials
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

	// Setting same site as strict for CSRF
	c.SetSameSite(http.SameSiteStrictMode)

	// TODO: change secure to true when using FE instead of postman
	// Set the cookie to mark user as logged in
	c.SetCookie(
		"sessionID",
		sessionID,
		ttlSec,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(status, bson.M{
		"message": "Logged in succesfully",
	})
}

func Welcome(c *gin.Context) {
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, bson.M{
			"error": "user unauthorized",
		})
		return
	}

	status, username, err := authController.IsAuthenticated(cookie)
	if err!=nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(status, gin.H{
		"message": fmt.Sprintf("User %s successfully logged in", username),
	})
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, bson.M{
			"error": "No logged in user found",
		})
		return
	}

	status, err := authController.Logout(cookie)

	if err!=nil {
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: change secure to true when using FE instead of postman
	c.SetCookie(
		"sessionID",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(status, gin.H{
		"message": "User logged out Successfully",
	})
}
