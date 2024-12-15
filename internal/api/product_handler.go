package api

import (
	"net/http"
	"strconv"
	"time"

	"product-management-system/internal/models"
	"product-management-system/internal/service"
	"product-management-system/internal/shared"
	"product-management-system/pkg/logger"
	"product-management-system/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(ps *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: ps,
	}
}

// CreateProduct handles the POST /products endpoint
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Measure request processing time
	start := time.Now()

	// Bind JSON input to product model
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Log.WithError(err).Error("Invalid product input")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Validate product input
	if err := utils.ValidateProduct(product); err != nil {
		logger.Log.WithError(err).Error("Product validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID, exists := c.Get("user_id")
	if !exists {
		logger.Log.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}
	product.UserID = userID.(uint)

	// Create product
	createdProduct, err := h.productService.CreateProduct(&product)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to create product")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Product creation failed",
			"details": err.Error(),
		})
		return
	}

	// Log request processing time
	duration := time.Since(start)
	logger.Log.WithFields(logrus.Fields{
		"product_id": createdProduct.ID,
		"duration":   duration,
	}).Info("Product created successfully")

	c.JSON(http.StatusCreated, createdProduct)
}

// GetProductByID handles the GET /products/:id endpoint
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	start := time.Now()

	// Parse product ID from URL
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.WithError(err).Error("Invalid product ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	// Retrieve product with caching
	product, err := h.productService.GetProductByID(uint(productID))
	if err != nil {
		if err == service.ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			logger.Log.WithError(err).Error("Failed to retrieve product")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Product retrieval failed",
				"details": err.Error(),
			})
		}
		return
	}

	// Log cache hit/miss
	duration := time.Since(start)
	logger.Log.WithFields(logrus.Fields{
		"product_id": product.ID,
		"duration":   duration,
	}).Info("Product retrieved")

	c.JSON(http.StatusOK, product)
}

// ListProducts handles the GET /products endpoint with filtering
func (h *ProductHandler) ListProducts(c *gin.Context) {
	start := time.Now()

	// Parse query parameters
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
	productName := c.Query("product_name")

	// Create filter struct
	filter := shared.ProductFilter{
		UserID:       uint(userID),
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		ProductName:  productName,
	}

	// Retrieve filtered products
	products, err := h.productService.ListProducts(filter)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to list products")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Product listing failed",
			"details": err.Error(),
		})
		return
	}

	// Log request processing
	duration := time.Since(start)
	logger.Log.WithFields(logrus.Fields{
		"products_count": len(products),
		"duration":       duration,
	}).Info("Products listed")

	c.JSON(http.StatusOK, products)
}


