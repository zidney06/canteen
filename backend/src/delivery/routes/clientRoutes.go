package routes

import (
	"server/src/delivery/handlers"
	"server/src/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func ClientRoutes(r *gin.RouterGroup, clientHandlers *handlers.ClientHandler) {
	// Inisialisasi handler

	// Grouping rute
	clientGroup := r.Group("/client")
	{
		clientGroup.GET("/", middlewares.AdminMiddleware(), clientHandlers.GetClientList)
		clientGroup.GET("/:clientId", middlewares.AdminMiddleware(), clientHandlers.GetSpecificClient)
		clientGroup.GET("/cs/:token", clientHandlers.GetClientSecret)
		clientGroup.POST("/", middlewares.AdminMiddleware(), clientHandlers.CreateNewClient)
		clientGroup.PATCH("/", middlewares.AdminMiddleware(), clientHandlers.CreateNewClientSecret)
	}
}
