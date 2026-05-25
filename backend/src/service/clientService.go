package service

import (
	"log"
	"server/src/repository"
	"server/src/types"
	"server/src/utils"

	"github.com/google/uuid"
)

type ClientService struct {
	Repo *repository.ClientRepo
}

func NewClientService(repo *repository.ClientRepo) *ClientService {
	return &ClientService{
		Repo: repo,
	}
}

func (s *ClientService) GetClientListService() ([]types.Client, error) {
	clients, err := s.Repo.GetClientList()
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			return nil, &types.ServiceError{
				Message:    customErr.Message,
				HttpStatus: 500,
			}
		}
	}

	return clients, nil
}

func (s *ClientService) GetSpecificClientService(clientId string) (types.Client, error) {
	client, err := s.Repo.GetClientById(clientId)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			return types.Client{}, &types.ServiceError{
				Message:    customErr.Message,
				HttpStatus: 500,
			}
		}
	}

	return client, nil
}

func (s *ClientService) GetClientSecretService(token string) (types.NewClientSecretType, error) {
	newClientSecret, err := s.Repo.GetClientSecretRepo(token)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			if customErr.Type == 0 {
				return types.NewClientSecretType{}, &types.ServiceError{
					Message:    customErr.Message,
					HttpStatus: 500,
				}
			} else {
				return types.NewClientSecretType{}, &types.ServiceError{
					Message:    customErr.Message,
					HttpStatus: 404,
				}
			}
		}
	}

	return newClientSecret, nil
}

func (s *ClientService) CreateNewClient(clientName string) (string, error) {
	var clientTamplate types.NewClientTamplate

	clientSecretOneTime := utils.GenerateClientSecret()

	clientTamplate.ID = uuid.New()
	clientTamplate.ClientName = clientName
	clientTamplate.ClientSecret = utils.HashingSecret(clientSecretOneTime)
	clientTamplate.ClientKey = utils.GenerateClientKey()

	newClient, err := s.Repo.CreateNewClient(clientTamplate)
	if err != nil {
		return "", &types.ServiceError{
			Message:    "Internal server error",
			HttpStatus: 500,
		}
	}

	token, err := utils.GenerateToken()
	if err != nil {
		log.Println(err)
		return "", &types.ServiceError{
			Message:    "Internal server error",
			HttpStatus: 500,
		}
	}

	// karena di sebeah kiri hanya ada satu variabel dan itu sudah pernah di deklarasikan sebelumnya
	// maka golang meganggap bahwa ini sedang memasukan data baru.
	// cara solvenya adalah dengna menggunakan = bukan :=
	err = s.Repo.SetToRedis(newClient.ClientKey, clientSecretOneTime, token)
	if err != nil {
		return "", &types.ServiceError{
			Message:    "Internal server error",
			HttpStatus: 500,
		}
	}

	return token, nil
}

func (s *ClientService) CreateNewClientSecret(clientKey string) (string, error) {
	newSecret := utils.GenerateClientSecret()
	hashedNewSecret := utils.HashingSecret(newSecret)

	err := s.Repo.SetNewClientSecret(clientKey, hashedNewSecret)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			if customErr.Type == 0 {
				return "", &types.ServiceError{
					Message:    customErr.Message,
					HttpStatus: 500,
				}
			} else {
				return "", &types.ServiceError{
					Message:    customErr.Message,
					HttpStatus: 404,
				}
			}
		}
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return "", &types.ServiceError{
			Message:    "Error while generating token!",
			HttpStatus: 500,
		}
	}

	err = s.Repo.SetToRedis(clientKey, newSecret, token)
	if err != nil {
		return "", &types.ServiceError{
			Message:    "Internal server error!",
			HttpStatus: 500,
		}
	}

	return token, nil
}
