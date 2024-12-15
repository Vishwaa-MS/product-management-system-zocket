package main

import (
	"fmt"
	"log"
	"product-management-system/config"
	"product-management-system/internal/api"
	"product-management-system/internal/cache"
	"product-management-system/internal/queue"
	"product-management-system/internal/repository"
	"product-management-system/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configurations
	cfg := config.LoadConfig()

	// Initialize database connection
	db := repository.InitPostgresDB(repository.DatabaseConfig{
	Host:     cfg.Database.Host,
	Port:     cfg.Database.Port,
	User:     cfg.Database.User,
	Password: cfg.Database.Password,
	DBName:   cfg.Database.DBName,
})

	// Initialize Redis cache
	redisCache := cache.NewRedisCache(cache.CacheConfig{
	Host:     cfg.Redis.Host,
	Port:     cfg.Redis.Port,
	Password: cfg.Redis.Password,
})

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)

	rabbitMQ := queue.NewRabbitMQ(cfg.RabbitMQ.Host, cfg.RabbitMQ.Port, cfg.RabbitMQ.QueueName)

	// Initialize services
	productService := service.NewProductService(*productRepo, *redisCache)
	imageProcessor := service.NewImageProcessor(rabbitMQ)

	// Start image processing queue consumer
	go imageProcessor.ConsumeImageProcessingQueue()

	// Setup Gin router
	router := gin.Default()
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize product handler
	productHandler := api.NewProductHandler(productService)

	// Define routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", productHandler.CreateProduct)
		v1.GET("/products/:id", productHandler.GetProductByID)
		v1.GET("/products", productHandler.ListProducts)
	}

	// Start server
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
