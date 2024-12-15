package service

import (
	"product-management-system/internal/models"
	"product-management-system/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	Repo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(user *models.User) error {
	return s.Repo.CreateUser(user)
}

// GetUserByEmail retrieves a user by their email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(email)
}
