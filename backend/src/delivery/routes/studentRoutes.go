package routes

import (
	"server/src/delivery/handlers"

	"github.com/gin-gonic/gin"
)

func Student(r *gin.RouterGroup, handlers *handlers.StudentHandler) {
	// Grouping rute
	studentGroup := r.Group("/student")
	{
		studentGroup.POST("/qr", handlers.CreateQrCode)
	}
}
