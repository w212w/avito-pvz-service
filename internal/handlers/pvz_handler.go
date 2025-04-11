package handlers

import (
	"avito-pvz-service/internal/middleware"
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type PVZHandler struct {
	pvzService *services.PVZService
}

func NewPVZHandler(pvzService *services.PVZService) *PVZHandler {
	return &PVZHandler{pvzService: pvzService}
}

func (h *PVZHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var pvzReq models.PVZ
	err := json.NewDecoder(r.Body).Decode(&pvzReq)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	role, _ := ctx.Value(middleware.RoleKey).(string)
	if role != services.RoleModerator {
		writeJSONError(w, http.StatusForbidden, "access denied")
		return
	}

	newPVZ, err := h.pvzService.CreatePVZ(pvzReq.City)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Error creating PVZ: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newPVZ); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode pvz")
	}
}
