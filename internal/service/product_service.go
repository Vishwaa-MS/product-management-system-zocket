package service

import (
	"errors"
	"product-management-system/internal/cache"
	"product-management-system/internal/models"
	"product-management-system/internal/repository"
	"product-management-system/internal/shared"
)


var ErrProductNotFound = errors.New("product not found")

// ProductService handles business logic for products
type ProductService struct {
	Repo repository.ProductRepository
	Cache cache.RedisCache
}

// ProductFilter represents filtering criteria for listing products


// NewProductService creates a new ProductService
func NewProductService(repo repository.ProductRepository, cache cache.RedisCache) *ProductService {
	return &ProductService{Repo: repo,Cache: cache}
}

// CreateProduct adds a new product
func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	err := s.Repo.CreateProduct(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}


// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	return s.Repo.GetProductByID(id)
}

// ListProducts retrieves all products for a user with optional filters
func (s *ProductService) ListProducts(filter shared.ProductFilter) ([]models.Product, error) {
	return s.Repo.ListProducts(filter)
}


