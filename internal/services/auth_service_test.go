package services_test

import (
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) CreateUser(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestRegister_ValidData(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	email := "user@example.com"
	password := "password123"
	role := services.RoleEmployee

	userRepo.On("CreateUser", mock.Anything).Return(nil)

	user, err := service.Register(email, password, role)

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
	userRepo.AssertExpectations(t)
}

func TestRegister_InvalidRole(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	_, err := service.Register("user@example.com", "password", "invalid")
	assert.EqualError(t, err, "invalid role")
}

func TestLogin_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	email := "user@example.com"
	password := "password123"
	hashedPassword, _ := services.HashPassword(password)

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		Role:         services.RoleModerator,
		PasswordHash: hashedPassword,
	}

	userRepo.On("GetUserByEmail", email).Return(user, nil)
	token, err := service.Login(email, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestLogin_InvalidPassword(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	hashedPassword, _ := services.HashPassword("rightpass")
	user := &models.User{
		ID:           uuid.New(),
		Email:        "user@example.com",
		PasswordHash: hashedPassword,
	}

	userRepo.On("GetUserByEmail", "user@example.com").Return(user, nil)

	_, err := service.Login("user@example.com", "wrongpass")
	assert.EqualError(t, err, "invalid password")
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	userRepo.On("GetUserByEmail", "missing@example.com").Return((*models.User)(nil), nil)

	_, err := service.Login("missing@example.com", "pass")
	assert.EqualError(t, err, "user not found")
}

func TestDummyLogin_Valid(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	token, err := service.DummyLogin(services.RoleEmployee)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	userRepo := new(mockUserRepo)
	service := services.NewAuthService(userRepo, "secret")

	_, err := service.DummyLogin("admin")
	assert.EqualError(t, err, "invalid role")
}
