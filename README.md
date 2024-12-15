# Product Management System

Welcome to the Product Management System! This project is a backend application built using Golang, designed to manage products with a focus on asynchronous processing, caching, logging, and high performance.

## Overview

The **Product Management System** is a scalable and modular backend application designed to manage product-related data, with functionalities such as CRUD operations for products, caching with Redis, and integration with a message queue for asynchronous processing. The system is built with Go and follows a clean architecture to ensure maintainability, scalability, and ease of understanding.
## Table of Contents

- [Project Architecture](#project-architecture)
- [Features](#features)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Testing](#testing)
- [Performance Optimization](#performance-optimization)
- [License](#license)

## Project Architecture

```
product-management-system/
├── cmd/
│   └── main.go               # Entry point of the application
├── internal/
│   ├── api/                  # API route handlers
│   │   ├── product_handler.go
│   │   └── middleware.go
│   ├── models/               # Database models
│   │   ├── product.go
│   │   └── user.go
│   ├── repository/           # Database interaction layer
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   ├── service/              # Business logic
│   │   ├── product_service.go
│   │   └── image_processor.go
│   ├── cache/                # Redis caching implementation
│   │   └── redis_cache.go
│   ├── queue/                # Message queue handling
│   │   └── rabbitmq.go
├── configs/                  # Configuration files
│   ├── config.yaml
│   └── database.sql
├── pkg/                      # Shared packages
│   ├── logger/               # Logging utility
│   │   └── logger.go
│   └── utils/                # Utility functions
│       └── validators.go
├── go.mod
└── README.md
```

The project follows a clean, modular architecture based on separation of concerns, implementing the following design patterns:

- **Layered Architecture**
- **Repository Pattern**
- **Dependency Injection**
- **Middleware-based Request Processing**

### Key Architecture Components

1. **Command Layer (`cmd/`)**:

   - Application entry point
   - Initializes and wires up all system components
   - Handles application startup and configuration loading

2. **Internal Layer (`internal/`)**:

   - Core business logic and implementation
   - Strictly private packages not importable by external projects

   #### Subcomponents:

   - **API Handlers**: Route management and request processing
   - **Models**: Database schema and data structures
   - **Repository**: Database interaction and data access
   - **Service**: Business logic and domain-specific operations
   - **Cache**: Redis-based caching mechanism
   - **Queue**: RabbitMQ message queue integration

3. **Package Layer (`pkg/`)**:

   - Shared, reusable packages
   - Utility functions and cross-cutting concerns
   - Can be potentially used by other projects

4. **Configuration Layer (`configs/`)**:
   - Application and infrastructure configurations
   - Database schema definitions

## Features

- **Product Management**

  - CRUD operations for products
  - Advanced filtering and search
  - Image processing for product images
  - Bulk import/export capabilities

- **Caching**

  - Redis-based intelligent caching
  - Automatic cache invalidation
  - Performance optimization

- **Messaging**

  - RabbitMQ integration for asynchronous processing
  - Event-driven architecture
  - Background job management

- **Security**
  - Middleware-based authentication
  - Request validation
  - Role-based access control

## Getting Started

### Prerequisites

- Go 1.20+
- Redis
- RabbitMQ
- PostgreSQL

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/product-management-system.git
cd product-management-system
```

2. Install dependencies

```bash
go mod download
```

3. Configure environment

```bash
cp configs/config.yaml.example configs/config.yaml
# Edit configuration as needed
```

4. Run database migrations

```bash
go run cmd/migrate/main.go
```

5. Start the application

```bash
go run cmd/main.go
```

## Project Structure

### Directory Breakdown

- `cmd/`: Application entry points
- `internal/`:
  - `api/`: HTTP route handlers and middleware
  - `models/`: Database models and structures
  - `repository/`: Database interaction logic
  - `service/`: Business logic implementations
  - `cache/`: Caching mechanisms
  - `queue/`: Message queue handling
- `configs/`: Configuration files
- `pkg/`:
  - `logger/`: Logging utilities
  - `utils/`: Shared utility functions

## Configuration

Configuration is managed via `configs/config.yaml`:

```yaml
server:
  port: 8080
  debug: false

database:
  host: localhost
  port: 5432
  name: productdb
  user: dbuser
  password: secret

redis:
  host: localhost
  port: 6379
  password: "redis"

rabbitmq:
  url: amqp://guest:guest@localhost:5672/
```

## API Endpoints

### Products

- `GET /products`: List products
- `GET /products/{id}`: Get product details
- `POST /products`: Create new product

### Authentication

- `POST /login`: User authentication
- `POST /register`: User registration

### Asynchronous Image Processing

Product images are processed asynchronously to ensure non-blocking operations and enhance performance. Upon creating a product, image URLs are added to a RabbitMQ queue. The image processor service listens for these messages, compresses the images, and updates the database with compressed image URLs.

### Caching

Redis is used to cache product data to reduce database load and improve response times. The cache is invalidated whenever product data is updated to ensure real-time accuracy.

### Logging

Structured logging is implemented using Logrus. All requests, responses, and processing details are logged, including specific events in the image processing service.

### Error Handling

Robust error handling is implemented across all components. This includes retry mechanisms for asynchronous processing failures and dead-letter queues for unprocessable messages.

## Development

1. Set up the configuration file `config.yaml` and the database schema using the provided `database.sql` script.
2. Ensure that Redis and RabbitMQ are running if you're using the caching and message queue features.
3. Start the application:
    ```bash
    go run cmd/main.go
    ```

The application will be available at `http://localhost:8080`, and you can use Postman or any other HTTP client to interact with the API.

### Running Tests

```bash
go test ./...
```

### Running Specific Tests

```bash
go test ./internal/service
```

## Performance Optimization

- Implemented Redis caching for frequently accessed data
- RabbitMQ for asynchronous task processing
- Efficient database query optimization
- Prepared statement usage
- Connection pooling

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Vishwaa M - vishwaams.03@gmail.com

Project Link: [https://github.com/yourusername/product-management-system](https://github.com/yourusername/product-management-system)
