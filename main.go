package main

import (
	"log"
	"net/http"
	"os"

	routes "streamit/router"
	"streamit/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func main() {
	utils.LoadEnv()         // Load environment variables
	utils.SetupGoogleAuth() // Initialize Google OAuth

	r := gin.Default()

	// Middleware for adapting gothic with Gin
	r.Use(func(c *gin.Context) {
		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return "google", nil
		}
		c.Next()
	})

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))) // Initialize session store

	// Auth Routes
	auth := r.Group("/")
	{
		routes.AuthHandler(auth)
	}

	log.Fatal(r.Run(":8080"))
}
