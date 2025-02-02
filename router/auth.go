package routes

import (
	"streamit/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthHandler(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/google", controllers.GoogleLoginHandler)
	auth.Get("/google/callback", controllers.GoogleCallbackHandler)
	auth.Get("/", controllers.HomeHandler)
}
