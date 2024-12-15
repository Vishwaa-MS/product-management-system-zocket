package utils

import (
	"errors"
	"product-management-system/internal/models"
)

func ValidateProduct(product models.Product) error {
	if product.ProductName == "" {
		return errors.New("product name is required")
	}
	if product.ProductPrice <= 0 {
		return errors.New("product price must be positive")
	}
	return nil
}

func ValidateProductUpdate(product models.Product) error {
	// Similar validation, but can be less strict
	return nil
}