package handlers

import (
	"avito-pvz-service/internal/middleware"
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
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

	role, _ := ctx.Value(middleware.RoleKey).(string)
	if role != services.RoleModerator {
		writeJSONError(w, http.StatusForbidden, "access denied")
		return
	}

	var pvzReq models.PVZ
	err := json.NewDecoder(r.Body).Decode(&pvzReq)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
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

func (h *PVZHandler) GetPVZList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	role, _ := ctx.Value(middleware.RoleKey).(string)
	if role != services.RoleModerator && role != services.RoleEmployee {
		writeJSONError(w, http.StatusForbidden, "access denied")
		return
	}

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var (
		startDate *time.Time
		endDate   *time.Time
	)

	if startDateStr != "" {
		startDateParsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid startDate")
			return
		}
		startDate = &startDateParsed
	}

	if endDateStr != "" {
		endDateParsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid endDate")
			return
		}
		endDate = &endDateParsed
	}

	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			writeJSONError(w, http.StatusBadRequest, "invalid page")
			return
		}
		page = p
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 30 {
			writeJSONError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = l
	}

	pvzList, err := h.pvzService.GetPVZList(startDate, endDate, page, limit)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(pvzList); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

}

func (h *PVZHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	role, _ := ctx.Value(middleware.RoleKey).(string)
	if role != services.RoleModerator && role != services.RoleEmployee {
		writeJSONError(w, http.StatusForbidden, "access denied")
		return
	}

	var req struct {
		PVZID uuid.UUID `json:"pvzId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	reception, err := h.pvzService.CreateReception(req.PVZID)

	if err != nil {
		if err == services.ErrReceptionExists {
			writeJSONError(w, http.StatusBadRequest, "reception already in progress")
			return
		}
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("Error creating reception: %s", err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reception); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode reception")
	}
}
