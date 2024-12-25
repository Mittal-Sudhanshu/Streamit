package controllers

import (
	"net/http"
	"streamit/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func GoogleLoginHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// GoogleCallbackHandler handles the callback after Google login
func GoogleCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate a JWT token
	token, error := utils.GenerateJWT(user.UserID)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	// Get the session from the store
	session, err := gothic.Store.Get(c.Request, "auth-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Store the JWT token in the session
	session.Values["jwt"] = token
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"token":   token,
		"user": gin.H{
			"name":  user.FirstName + " " + user.LastName,
			"email": user.Email,
			"photo": user.AvatarURL,
		},
	})
}
