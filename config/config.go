package config

import (
	"log"
	"os"
	"product-management-system/internal/repository"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the config file
type Config struct {
	Server struct {
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"`
	} `yaml:"server"`
	Database repository.DatabaseConfig `yaml:"database"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	RabbitMQ struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
	S3 struct {
		Bucket string `yaml:"bucket"`
		Region string `yaml:"region"`
	} `yaml:"s3"`
}

// LoadConfig loads configuration from config.yaml
func LoadConfig() *Config {
	file, err := os.Open("configs/config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	return &cfg
}
