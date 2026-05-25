package routes

import (
	"server/src/delivery/handlers"
	"server/src/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup, transactionHandlers *handlers.TransactionHandler) {

	transactionGroup := r.Group("/transaction")
	{
		transactionGroup.POST("/", middlewares.JWTMiddleware(), transactionHandlers.MakeTransaction)
	}
}
