package service

import (
	"errors"
	"server/src/repository"
	"server/src/types"
)

type ProductService struct {
	Repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}

func (s *ProductService) GetProductList(clientId string) ([]types.Product, error) {
	if clientId == "" {
		return nil, &types.ServiceError{
			Message:    "Client ID is required!",
			HttpStatus: 400,
		}
	}

	products, err := s.Repo.GetProductList(clientId)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			return nil, &types.ServiceError{
				Message:    customErr.Message,
				HttpStatus: 500,
			}
		}
	}

	return products, nil
}

func (s *ProductService) AddProducts(items []types.Product, clientId string) (int, error) {
	if len(items) == 0 {
		return 0, &types.ServiceError{
			Message:    "Items are required!",
			HttpStatus: 400,
		}
	} else if clientId == "" {
		return 0, &types.ServiceError{
			Message:    "Client ID is required!",
			HttpStatus: 400,
		}
	}

	rowsAffected, err := s.Repo.AddProducts(items, clientId)
	if err != nil {
		if customErr, ok := err.(*types.RepoError); ok {
			return 0, &types.ServiceError{
				Message:    customErr.Message,
				HttpStatus: 500,
			}
		}
		return 0, err
	}

	return rowsAffected, nil
}

func (s *ProductService) UpdateProduct(productId string, itemName string, price float64) (types.Product, error) {
	var repoErr *types.RepoError
	if productId == "" {
		return types.Product{}, &types.ServiceError{
			Message:    "Product ID is required!",
			HttpStatus: 400,
		}
	} else if itemName == "" {
		return types.Product{}, &types.ServiceError{
			Message:    "Item name is required!",
			HttpStatus: 400,
		}
	} else if price <= 0 {
		return types.Product{}, &types.ServiceError{
			Message:    "Price must be greater than 0!",
			HttpStatus: 400,
		}
	}

	updatedProduct, err := s.Repo.UpdateProduct(productId, itemName, price)
	if err != nil && errors.As(err, &repoErr) {
		switch repoErr.Type {
		case 1:
			return types.Product{}, &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 404,
			}
		default:
			return types.Product{}, &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 500,
			}
		}
	}

	return updatedProduct, nil
}

func (s *ProductService) DeleteProduct(productIds []string) error {
	var repoErr *types.RepoError
	if len(productIds) == 0 {
		return &types.ServiceError{
			Message:    "array of productIds must not zero!",
			HttpStatus: 400,
		}
	}

	err := s.Repo.DeleteProduct(productIds)
	if err != nil && errors.As(err, &repoErr) {
		switch repoErr.Type {
		case 2:
			return &types.ServiceError{
				Message:    repoErr.Message,
				HttpStatus: 400,
			}
		default:
			return &types.ServiceError{
				Message:    "Internal server error!",
				HttpStatus: 500,
			}
		}
	}

	return nil
}
