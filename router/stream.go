package routes

import (
	"streamit/controllers"
	"streamit/middleware"

	"github.com/gofiber/fiber/v2"
)

func StreamRouter(router fiber.Router) {
	// StreamRouter is a function that returns a gin router for the stream endpoints
	stream := router.Group("/stream")
	router.Use(middleware.AuthMiddleware)
	stream.Post("/newStream", controllers.CreateNewStreamKey)
	stream.Post("/channel", controllers.CreateNewChannel)
}
