package handlers

import (
	"net/http"
	"server/src/service"
	"server/src/types"

	"github.com/gin-gonic/gin"
)

/*
 * alur pembuatan client
 * 1. server menerima client name
 * 2. server membuat client key dan client secret
 * 3. data disimpan ke database
 * 4. buat one time url yang akan admin kirimkan ke client yan berisi client key dan client secret
 * 7. saat client mau login, client memasukan cs dan ck
 * 8. server validasi
 * 9. kalau lolos, kirim jwt ke client + client keynya didalam isi jwt
 *
 * PEMBUATAN CLIENT SECRET BARU
 * kita memakai one time url. One-Time URL (atau sering disebut Burning Link) adalah sebuah mekanisme keamanan di mana sebuah tautan (URL)
 * hanya dapat diakses tepat satu kali. Setelah diklik atau dibuka,
 * data yang dikaitkan dengan URL tersebut akan langsung dihapus dari server, sehingga link tersebut tidak bisa digunakan lagi.
 */

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

// GET

// POST
func (h *AuthHandler) ClientLogin(c *gin.Context) {
	var body types.ClientLoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// service
	tokenString, err := h.Service.ClientLoginService(body.ClientKey, body.ClientSecret)
	if err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Client login successful",
		"jwtToken": tokenString,
	})
}

func (h *AuthHandler) ClientLogout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if err := h.Service.ClientLogoutService(tokenString); err != nil {
		if customErr, ok := err.(*types.ServiceError); ok {
			c.JSON(customErr.HttpStatus, gin.H{
				"message": customErr.Message,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Client logout successfully",
	})
}
