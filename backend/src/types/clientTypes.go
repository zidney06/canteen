package types

import (
	"github.com/google/uuid"
)

type Client struct {
	ID         uuid.UUID `json:"id"`
	ClientName string    `json:"clientName"`
	ClientKey  string    `json:"clientKey"`
}

type CreateNewClientSecretBody struct {
	ClientKey string `json:"clientKey" binding:"required"`
}

type NewClientSecretType struct {
	ClientKey    string
	ClientSecret string
}

type CreateNewClientBody struct {
	ClientName string `json:"clientName" binding:"required"`
}

type NewClientTamplate struct {
	ID           uuid.UUID
	ClientName   string
	ClientKey    string
	ClientSecret string
}
