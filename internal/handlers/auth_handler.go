package handlers

import (
	"avito-pvz-service/internal/services"
	logger "avito-pvz-service/pkg"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type registerResponse struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("handler /register started")
	w.Header().Set("Content-Type", "application/json")
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := registerResponse{
		Email: user.Email,
		Role:  user.Role,
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("handler /login started")
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	resp := tokenResponse{Token: token}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode response")
	}
}

type dummyLoginRequest struct {
	Role string `json:"role"`
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("handler /dummyLogin started")
	w.Header().Set("Content-Type", "application/json")
	var req dummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Role != services.RoleEmployee && req.Role != services.RoleModerator {
		writeJSONError(w, http.StatusBadRequest, "invalid role")
		return
	}

	token, err := h.authService.DummyLogin(req.Role)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to generate token")
		return
	}

	resp := tokenResponse{Token: token}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode response")
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Message: message})
}
