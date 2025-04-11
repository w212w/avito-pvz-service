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

type AuthService struct {
	userRepo  repository.UserRepoInterface
	secretKey string
}

func NewAuthService(userRepo repository.UserRepoInterface, secretKey string) *AuthService {
	return &AuthService{userRepo: userRepo, secretKey: secretKey}
}

func (s *AuthService) Register(email, password, role string) (*models.User, error) {
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

func (s *AuthService) Login(email, password string) (string, error) {
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

	return s.generateToken(user.ID, user.Role)
}

func (s *AuthService) DummyLogin(role string) (string, error) {
	if role != RoleEmployee && role != RoleModerator {
		return "", errors.New("invalid role")
	}

	dummyID := uuid.New()
	return s.generateToken(dummyID, role)
}

func (s *AuthService) generateToken(userID uuid.UUID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ParseToken(tokenString string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, "", errors.New("invalid user_id format")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, "", errors.New("user_id is not a valid UUID")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return uuid.Nil, "", errors.New("invalid role format")
		}

		return userID, role, nil
	}

	return uuid.Nil, "", errors.New("invalid token")
}
