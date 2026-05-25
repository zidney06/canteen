package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"server/src/types"

	"github.com/google/uuid"
)

// GenerateClientKey menggunakan UUID dengan prefix agar mudah dikenali
func GenerateClientKey() string {
	return fmt.Sprintf("CK_%s", uuid.New().String())
}

// GenerateClientSecret menggunakan random bytes (lebih aman untuk secret)
func GenerateClientSecret() string {
	b := make([]byte, 32) // 32 bytes = 256 bit
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func HashingSecret(secret string) string {
	salt := os.Getenv("HASHING_SALT")

	hashResult := sha256.Sum256([]byte(secret + salt))
	return hex.EncodeToString(hashResult[:])
}

func PrettyJson(input any) {
	prettyJSON, _ := json.MarshalIndent(input, "", "    ")
	log.Println("Pretty result: ", string(prettyJSON))
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RepoToServiceError(err error) *types.ServiceError {
	var repoErr *types.RepoError

	if errors.As(err, &repoErr) {
		return repoErr.ToServiceError()
	}

	return &types.ServiceError{
		Message:    "Internal server error.",
		HttpStatus: 500,
	}
}

// 1. Fungsi untuk mengenkripsi string plaintext
func Encrypt(plaintext string, key []byte) (string, error) {
	// Membuat block cipher baru menggunakan AES dan key kita
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Menggunakan mode GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Membuat Nonce (Number used once) acak.
	// Nonce WAJIB unik setiap kali melakukan enkripsi dengan key yang sama.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Melakukan enkripsi. Hasilnya menggabungkan nonce di depan data terenkripsi.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Mengubah hasil biner menjadi string Hex agar mudah disimpan/dikirim
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// 2. Fungsi untuk mendekripsi kembali string ciphertextHex
func Decrypt(ciphertextHex string, key []byte) (string, error) {
	// Mengubah kembali string Hex menjadi bytes biner
	ciphertext, err := base64.URLEncoding.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext terlalu pendek")
	}

	// Memisahkan nonce dan ciphertext asli
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Dekripsi dan validasi integritas data
	plaintextBytes, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gagal mendekripsi (key salah atau data korup): %v", err)
	}

	return string(plaintextBytes), nil
}
