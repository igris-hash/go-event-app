package service

import (
	"errors"
	"time"

	"github.com/igris-hash/go-event-app/internal/model"
	"github.com/igris-hash/go-event-app/internal/repository"
	"github.com/igris-hash/go-event-app/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register creates a new user
func (s *UserService) Register(input model.UserRegistration) (*model.User, error) {
	// Check if user already exists
	exists, err := s.repo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save user
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user
func (s *UserService) Login(input model.UserLogin) (string, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(input.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

// Update updates user information
func (s *UserService) Update(id int, input model.UserUpdate) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
