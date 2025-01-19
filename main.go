package main

import (
	"net/http"
	"os"
	"sync"

	routes "streamit/router"
	"streamit/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
)

func Init() {
	utils.LoadEnv()
	utils.SetupGoogleAuth()
	utils.ConnectDB()
}

func main() {
	Init()

	// WaitGroup to run Gin and RTMP servers concurrently
	var wg sync.WaitGroup

	// Start RTMP Server
	server := CreateRTMPServer()
	if server == nil {
		log.Fatal("Failed to create RTMP server")
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Starting RTMP server...")
			server.Start() // Assuming this blocks while the server is running
		}()
	}

	// Start Gin Server
	r := gin.Default()

	// Middleware for adapting gothic with Gin
	r.Use(func(c *gin.Context) {
		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return "google", nil
		}
		c.Next()
	})

	// Initialize session store
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	// Serve HLS files
	// r.Static("hls", "./hls")

	// Auth Routes
	auth := r.Group("/")
	routes.AuthHandler(auth)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting Gin server on port 8080...")
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Gin server failed: %v", err)
		}
	}()

	wg.Wait()
}
