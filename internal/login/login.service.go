// internal/login/service.go
package login

import (
	"github.com/brain-flowing-company/pprp-backend/apperror"
	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	AuthenticateUser(email, password string) (string, *apperror.AppError)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo,
	}
}

func (s *serviceImpl) AuthenticateUser(email, password string) (string, *apperror.AppError) {
	// Retrieve user by email
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", apperror.UserNotFound
	}

	// Check password
	if err := s.checkPassword(user, password); err != nil {
		return "", apperror.InvalidCredentials
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	// Add other claims as needed

	// Sign the token with a secret key
	secretKey := []byte("secret") // Replace with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", apperror.InternalServerError
	}

	return tokenString, nil
}

func (s *serviceImpl) checkPassword(user *models.User, password string) *apperror.AppError {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return apperror.InvalidCredentials
	}
	return nil
}
