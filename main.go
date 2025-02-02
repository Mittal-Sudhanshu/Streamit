package main

import (
	"encoding/json"
	"sync"

	routes "streamit/router"
	"streamit/utils"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(recover.New())
	api := app.Group("/")
	routes.AuthHandler(api)
	routes.StreamRouter(api)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting Fiber server on port 8080...")
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Gin server failed: %v", err)
		}
	}()

	wg.Wait()
}
