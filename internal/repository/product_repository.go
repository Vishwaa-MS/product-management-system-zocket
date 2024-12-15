package repository

import (
	"product-management-system/internal/models"
	"product-management-system/internal/shared"

	"gorm.io/gorm"
)

// ProductRepository handles database interactions for products
type ProductRepository struct {
	DB *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

// CreateProduct inserts a new product into the database

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	return r.DB.Create(product).Error
}


// GetProductByID retrieves a product by its ID
func (r *ProductRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.DB.First(&product, id).Error
	return &product, err
}

// ListProducts retrieves all products with optional filters
// ListProducts retrieves all products with optional filters
func (r *ProductRepository) ListProducts(filter shared.ProductFilter) ([]models.Product, error) {
	var products []models.Product
	query := r.DB.Model(&models.Product{}).Where("user_id = ?", filter.UserID)

	// Apply additional filters
	if filter.MinPrice > 0 {
		query = query.Where("product_price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("product_price <= ?", filter.MaxPrice)
	}
	if filter.ProductName != "" {
		query = query.Where("product_name ILIKE ?", "%"+filter.ProductName+"%")
	}

	// Execute the query
	err := query.Find(&products).Error
	return products, err
}

