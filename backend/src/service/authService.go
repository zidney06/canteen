package service

import (
	"log"
	"os"
	"server/src/repository"
	"server/src/types"
	"server/src/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	Repo *repository.AuthRepo
}

func NewAuthService(repo *repository.AuthRepo) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (s *AuthService) ClientLoginService(clientKey string, clientSecret string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	client, err := s.Repo.GetClientWithKey(clientKey)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			if customErr.Type == 2 {
				return "", &types.ServiceError{
					Message:    customErr.Message,
					HttpStatus: 404,
				}
			} else {
				return "", &types.ServiceError{
					Message:    "Internal server error!",
					HttpStatus: 500,
				}
			}
		}
	}

	hashedSecret := utils.HashingSecret(clientSecret)

	if hashedSecret != client.ClientSecret {
		return "", &types.ServiceError{
			Message:    "Invalid client secret",
			HttpStatus: 400,
		}
	}

	// kasih jwt selama 1 hari
	claims := types.MyClaims{
		ID:         client.ID,
		ClientName: client.ClientName,
		ClientKey:  client.ClientKey,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("%+v\n", err)
		return "", &types.ServiceError{
			Message:    "Internal server error",
			HttpStatus: 500,
		}
	}

	return tokenString, nil
}

func (s *AuthService) ClientLogoutService(tokenString string) error {
	var jwtToken string
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	if tokenString == "" {
		return &types.ServiceError{
			Message:    "Token required!",
			HttpStatus: 401,
		}
	} else if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		jwtToken = tokenString[7:] // Ambil setelah karakter ke-7
	}

	// Parse token
	token, err := jwt.ParseWithClaims(jwtToken, &types.MyClaims{}, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return &types.ServiceError{
			Message:    "Token invalid",
			HttpStatus: 401,
		}
	}

	if claims, ok := token.Claims.(*types.MyClaims); ok {

		// cek apakah token sudah ada di redis atau belum?
		isExist, msg := s.Repo.IsTokenBlacklisted(*claims)
		if isExist {
			return &types.ServiceError{
				Message:    msg,
				HttpStatus: 400,
			}
		}

		if err := s.Repo.SetBlacklistedToken(*claims); err != nil {
			return &types.ServiceError{
				Message:    err.Error(),
				HttpStatus: 500,
			}
		}
	}

	return nil
}
