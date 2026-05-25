package routes

import (
	"server/src/delivery/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	// Grouping rute
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.ClientLogin)
		authGroup.POST("/logout", authHandler.ClientLogout)
	}
}
