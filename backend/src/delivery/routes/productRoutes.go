package routes

import (
	"server/src/delivery/handlers"
	"server/src/delivery/middlewares"

	"github.com/gin-gonic/gin"
)

func Product(r *gin.RouterGroup, productHandlers *handlers.ProductHandler) {

	// Grouping rute
	productGroup := r.Group("/product")
	{
		productGroup.GET("/", middlewares.JWTMiddleware(), productHandlers.GetProductList)
		productGroup.POST("/", middlewares.JWTMiddleware(), productHandlers.AddProduct)
		productGroup.PATCH("/:productId", middlewares.JWTMiddleware(), productHandlers.UpdateProduct)
		productGroup.DELETE("/", middlewares.JWTMiddleware(), productHandlers.DeleteProducts)
	}
}
