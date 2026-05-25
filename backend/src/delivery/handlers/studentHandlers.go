package handlers

import (
	"io"
	"net/http"
	"server/src/service"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	Service *service.StudentService
}

func NewStudentHandler(serv *service.StudentService) *StudentHandler {
	return &StudentHandler{
		Service: serv,
	}
}

func (h *StudentHandler) CreateQrCode(c *gin.Context) {
	data := c.Query("text")
	if data == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter 'text' wajib diisi"})
		return
	}

	// 1. Panggil service untuk mendapatkan stream
	stream, err := h.Service.CreateQrCode(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Set Header karena format dari API luar adalah SVG
	c.Header("Content-Type", "image/svg+xml")

	// 3. Alirkan data langsung ke client menggunakan c.Stream
	c.Stream(func(w io.Writer) bool {
		// Pastikan stream ditutup setelah fungsi c.Stream selesai mengeksekusi data
		defer stream.Close()

		// Salurkan potongan biner QR Code langsung ke client
		_, err := io.Copy(w, stream)
		if err != nil {
			return false
		}
		return false // Berhenti karena data sudah habis dialirkan
	})
}
