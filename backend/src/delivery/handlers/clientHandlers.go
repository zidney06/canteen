package handlers

import (
	"net/http"

	"server/src/service"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	Service *service.ClientService
}

func NewClientHandler(service *service.ClientService) *ClientHandler {
	return &ClientHandler{
		Service: service,
	}
}

// GET
func (h *ClientHandler) GetClientList(c *gin.Context) {
	clients, err := h.Service.GetClientListService()
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get client list",
		"data":    clients,
	})
}

func (h *ClientHandler) GetSpecificClient(c *gin.Context) {
	clientId := c.Param("clientId")

	client, err := h.Service.GetSpecificClientService(clientId)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get client",
		"data":    client,
	})
}

func (h *ClientHandler) GetClientSecret(c *gin.Context) {
	token := c.Param("token")

	newClientSecret, err := h.Service.GetClientSecretService(token)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "This is for one time url",
		"data":    newClientSecret,
	})
}

// POST
func (h *ClientHandler) CreateNewClient(c *gin.Context) {
	var body types.CreateNewClientBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.CreateNewClient(body.ClientName)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Client created successfully",
		"token":   token,
	})
}

func (h *ClientHandler) CreateNewClientSecret(c *gin.Context) {
	var body types.CreateNewClientSecretBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Service.CreateNewClientSecret(body.ClientKey)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "New client secret created", "token": token})
}
