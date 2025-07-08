package main

import (
	"encoding/json"
	"os"
	"sync"

	routes "streamit/router"
	"streamit/utils"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/contrib/swagger"
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

	// Serve Swagger UI at /swagger
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}))
	// Serve the static swagger.json as /swagger/doc.json
	app.Static("/swagger/doc.json", "./swagger.json")
	api := app.Group("/")
	routes.AuthHandler(api)
	routes.StreamRouter(api)
	wg.Add(1)
	err := utils.InitS3(os.Getenv("AWS_S3_BUCKET"))
	if err != nil {
		log.Fatal("Failed to init S3:", err)
	}

	go func() {
		defer wg.Done()
		log.Info("Starting Fiber server on port 8080...")
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("Gin server failed: %v", err)
		}
	}()

	wg.Wait()
}
