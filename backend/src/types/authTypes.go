package types

import (
	"github.com/golang-jwt/jwt/v5"
)

// Request body type
type ClientLoginBody struct {
	ClientKey    string `json:"clientKey" binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
}

// Response type
type ClientAuth struct {
	ID           string `json:"id"`
	ClientKey    string `json:"clientKey"`
	ClientName   string `json:"clientName"`
	ClientSecret string `json:"clientSecret"`
}

// jwt
type MyClaims struct {
	ID         string `json:"id"`
	ClientName string `json:"username"`
	ClientKey  string `json:"clientKey"`
	jwt.RegisteredClaims
}
