package routes

import (
	"streamit/controllers"

	"github.com/gin-gonic/gin"
)

func AuthHandler(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.GET("/google", controllers.GoogleLoginHandler)
		auth.GET("/google/callback", controllers.GoogleCallbackHandler)
		auth.GET("/", controllers.HomeHandler)
	}
}
