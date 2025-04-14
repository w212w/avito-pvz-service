package handlers

import (
	pr "avito-pvz-service/internal/metrics"
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	logger "avito-pvz-service/pkg"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PVZHandler struct {
	pvzService services.PVZService
}

func NewPVZHandler(pvzService services.PVZService) *PVZHandler {
	return &PVZHandler{pvzService: pvzService}
}

func (h *PVZHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "201"
	defer func() {
		duration := time.Since(start)
		pr.LogRequestDuration(r.Method, duration)
		pr.LogRequestCount(r.Method, status)
	}()
	logger.Log.Info("handler /pvz POST started")
	w.Header().Set("Content-Type", "application/json")

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
	logger.Log.Info("handler /pvz GET started")
	w.Header().Set("Content-Type", "application/json")

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
	start := time.Now()
	status := "201"
	defer func() {
		duration := time.Since(start)
		pr.LogRequestDuration(r.Method, duration)
		pr.LogRequestCount(r.Method, status)
	}()
	logger.Log.Info("handler /receptions started")
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		PVZID uuid.UUID `json:"pvzId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PVZID == uuid.Nil {
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

func (h *PVZHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "201"
	defer func() {
		duration := time.Since(start)
		pr.LogRequestDuration(r.Method, duration)
		pr.LogRequestCount(r.Method, status)
	}()
	logger.Log.Info("handler /products started")
	w.Header().Set("Content-Type", "application/json")

	var req AddProductRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	product, err := h.pvzService.AddProduct(req.Type, req.PVZID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body or no reception")
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode product")
	}

}

type AddProductRequest struct {
	Type  string    `json:"type"`
	PVZID uuid.UUID `json:"pvzId"`
}

func (h *PVZHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("handler /pvz/{pvzId}/close_last_reception started")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	pvzIdStr := vars["pvzId"]
	if pvzIdStr == "" {
		writeJSONError(w, http.StatusBadRequest, "pvzId parameter is required")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	reception, err := h.pvzService.CloseLastReception(pvzId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request or reception closed")
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(reception); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to encode reception")
	}

}

func (h *PVZHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("handler /pvz/{pvzId}/delete_last_product started")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	pvzIdStr := vars["pvzId"]
	if pvzIdStr == "" {
		writeJSONError(w, http.StatusBadRequest, "pvzId parameter is required")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.pvzService.DeleteLastProduct(pvzId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}
	w.WriteHeader(http.StatusOK)

}
