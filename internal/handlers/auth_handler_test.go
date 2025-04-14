package handlers_test

import (
	"avito-pvz-service/internal/handlers"
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type mockAuthService struct {
	RegisterFunc      func(email, password, role string) (*models.User, error)
	LoginFunc         func(email, password string) (string, error)
	DummyLoginFunc    func(role string) (string, error)
	GenerateTokenFunc func(userID uuid.UUID, role string) (string, error)
}

func (m *mockAuthService) Register(email, password, role string) (*models.User, error) {
	return m.RegisterFunc(email, password, role)
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	return m.LoginFunc(email, password)
}

func (m *mockAuthService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	if m.GenerateTokenFunc == nil {
		return "", errors.New("error")
	}
	return m.GenerateTokenFunc(userID, role)
}

func (m *mockAuthService) DummyLogin(role string) (string, error) {
	return m.DummyLoginFunc(role)
}

func TestRegister_Success(t *testing.T) {
	authService := &mockAuthService{
		RegisterFunc: func(email, password, role string) (*models.User, error) {
			return &models.User{Email: email, Role: role}, nil
		},
	}
	h := handlers.NewAuthHandler(authService)

	reqBody := map[string]string{"email": "user@example.com", "password": "pass123", "role": "moderator"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestRegister_NoEmail(t *testing.T) {
	authService := &mockAuthService{
		RegisterFunc: func(email, password, role string) (*models.User, error) {
			return &models.User{Email: email, Role: role}, nil
		},
	}
	h := handlers.NewAuthHandler(authService)

	reqBody := map[string]string{"password": "pass123", "role": "moderator"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	authService := &mockAuthService{}
	h := handlers.NewAuthHandler(authService)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(`invalid json`)))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	authService := &mockAuthService{
		LoginFunc: func(email, password string) (string, error) {
			return "mock-token", nil
		},
	}
	h := handlers.NewAuthHandler(authService)

	body := map[string]string{"email": "user@example.com", "password": "pass123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestLogin_Unauthorized(t *testing.T) {
	authService := &mockAuthService{
		LoginFunc: func(email, password string) (string, error) {
			return "", errors.New("invalid credentials")
		},
	}
	h := handlers.NewAuthHandler(authService)

	body := map[string]string{"email": "wrong@example.com", "password": "badpass"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestDummyLogin_Success(t *testing.T) {
	authService := &mockAuthService{
		DummyLoginFunc: func(role string) (string, error) {
			return "dummy-token", nil
		},
	}
	h := handlers.NewAuthHandler(authService)

	body := map[string]string{"role": services.RoleModerator}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.DummyLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	authService := &mockAuthService{}
	h := handlers.NewAuthHandler(authService)

	body := map[string]string{"role": "invalid-role"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/dummy-login", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.DummyLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
