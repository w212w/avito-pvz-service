package services

import (
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

type AuthService interface {
	DummyLogin(role string) (string, error)
	Login(email, password string) (string, error)
	Register(email, password, role string) (*models.User, error)
	GenerateToken(userID uuid.UUID, role string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepoInterface
	secretKey string
}

func NewAuthService(userRepo repository.UserRepoInterface, secretKey string) AuthService {
	return &authService{userRepo: userRepo, secretKey: secretKey}
}

func (s *authService) Register(email, password, role string) (*models.User, error) {
	if role != RoleEmployee && role != RoleModerator {
		return nil, errors.New("invalid role")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Role:         role,
		PasswordHash: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (s *authService) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("empty email or password")
	}

	user, err := s.userRepo.GetUserByEmail(email)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	return s.GenerateToken(user.ID, user.Role)
}

func (s *authService) DummyLogin(role string) (string, error) {
	if role != RoleEmployee && role != RoleModerator {
		return "", errors.New("invalid role")
	}

	dummyID := uuid.New()
	return s.GenerateToken(dummyID, role)
}

func (s *authService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
