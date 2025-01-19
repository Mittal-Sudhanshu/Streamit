package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"streamit/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GoogleLoginHandler redirects to Google's OAuth login
func GoogleLoginHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func HomeHandler(c *gin.Context) {
	var users []bson.M

	cursor, err := utils.UserDb.Find(context.Background(), bson.M{}) // Use bson.M{} for an empty filter
	if err != nil {
		log.Fatalf("Error finding users: %v", err)
		return
	}
	defer cursor.Close(context.Background()) // Ensure the cursor is closed

	for cursor.Next(context.Background()) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			log.Fatalf("Error decoding user: %v", err)
			return
		}
		users = append(users, user)
	}

	// user, _ := utils.InitUsersCollection().Find(context.Background(), bson.M{})
	log.Print(err)
	// func(ctx *gin.Context) {
	c.JSON(200, gin.H{
		"users": users,
	})

	// c.JSON(http.StatusOK, gin.H{"message": "Welcome to StreamIt"})
}

// GoogleCallbackHandler handles the callback after Google login
func GoogleCallbackHandler(c *gin.Context) {
	// Complete user authentication with Google
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists in the database
	filter := bson.M{"google_id": user.UserID}

	var existingUser bson.M
	err = utils.UserDb.FindOne(context.Background(), filter).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user in database"})
		return
	}

	// If the user doesn't exist, insert them
	if err == mongo.ErrNoDocuments {
		newUser := bson.M{
			"google_id":      user.UserID,
			"name":           user.Name,
			"email":          user.Email,
			"picture":        user.AvatarURL,
			"verified_email": true,  // Assuming verified by Google
			"admin":          false, // Default admin flag
			"created_at":     time.Now(),
			"updated_at":     time.Now(),
		}

		_, err = utils.UserDb.InsertOne(context.Background(), newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
			return
		}
	}

	// Generate a JWT token
	token, jwtErr := utils.GenerateJWT(user.UserID)
	if jwtErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": jwtErr.Error()})
		return
	}

	// Store the JWT token in the session
	// session, err := gothic.Store.Get(c.Request, "auth-session")

	c.SetCookie("jwt", token, 60*60*24*7, "/", "localhost", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Respond to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   token,
		"user": gin.H{
			"name":  user.Name,
			"email": user.Email,
			"photo": user.AvatarURL,
		},
	})
}
