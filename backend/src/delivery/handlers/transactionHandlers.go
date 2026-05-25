package handlers

import (
	"log"
	"net/http"
	"server/src/delivery/middlewares"
	"server/src/service"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		Service: service,
	}
}

func (h *TransactionHandler) MakeTransaction(c *gin.Context) {
	var body types.MakeTransaactionBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := middlewares.GetCtxFromReq(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get client"})
		return
	}

	trs, err := h.Service.MakeTransaction(body.EncryptedId, body.Items, client.ClientID, client.ClientName)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "transaction made",
		"data":    trs,
	})
}
