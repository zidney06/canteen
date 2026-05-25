package handlers

import (
	"log"
	"net/http"
	"server/src/delivery/middlewares"
	"server/src/service"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		Service: service,
	}
}

// GET
func (h *ProductHandler) GetProductList(c *gin.Context) {
	client, err := middlewares.GetCtxFromReq(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get client"})
		return
	}

	products, err := h.Service.GetProductList(client.ClientID)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Succcess",
		"data":    products,
	})
}

// POST
// Penambahan produk baru harus berupa array agar frontend bisa langsung menmabahkan banyak produk sekaligus.
func (h *ProductHandler) AddProduct(c *gin.Context) {
	var body types.AddProductsBody
	client, err := middlewares.GetCtxFromReq(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get client"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := h.Service.AddProducts(body.Products, client.ClientID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "New products created successfully",
		"count":   rowsAffected,
	})
}

// PATCH
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productId := c.Param("productId")

	var body types.UpdateProductBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProduct, err := h.Service.UpdateProduct(productId, body.ItemName, body.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Success",
		"updatedProduct": updatedProduct,
	})
}

// DELETE
func (h *ProductHandler) DeleteProducts(c *gin.Context) {
	var body types.DeleteProductsBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.ProductsIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ProductIds value can't zero!",
		})
		return
	}

	err := h.Service.DeleteProduct(body.ProductsIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products deleted!",
	})
}
