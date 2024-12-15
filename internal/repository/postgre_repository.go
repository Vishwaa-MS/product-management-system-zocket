package repository

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// InitPostgresDB initializes and returns a GORM DB connection to PostgreSQL
func InitPostgresDB(cfg DatabaseConfig) *gorm.DB {
	// Build the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// Initialize the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Perform auto-migration for models if needed (optional)
	// db.AutoMigrate(&models.User{}, &models.Product{}) // Uncomment if models are available

	return db
}
