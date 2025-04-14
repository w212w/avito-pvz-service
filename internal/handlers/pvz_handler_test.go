package handlers_test

import (
	"avito-pvz-service/internal/handlers"
	"avito-pvz-service/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockPVZService struct {
	CreatePVZFunc          func(city string) (*models.PVZ, error)
	GetPVZListFunc         func(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
	CreateReceptionFunc    func(pvzID uuid.UUID) (*models.Reception, error)
	AddProductFunc         func(productType string, pvzID uuid.UUID) (*models.Product, error)
	CloseLastReceptionFunc func(pvzID uuid.UUID) (*models.Reception, error)
	DeleteLastProductFunc  func(pvzID uuid.UUID) error
}

func (m *mockPVZService) CreatePVZ(city string) (*models.PVZ, error) {
	return m.CreatePVZFunc(city)
}
func (m *mockPVZService) GetPVZList(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	return m.GetPVZListFunc(startDate, endDate, page, limit)
}
func (m *mockPVZService) CreateReception(pvzID uuid.UUID) (*models.Reception, error) {
	return m.CreateReceptionFunc(pvzID)
}
func (m *mockPVZService) AddProduct(productType string, pvzID uuid.UUID) (*models.Product, error) {
	return m.AddProductFunc(productType, pvzID)
}
func (m *mockPVZService) CloseLastReception(pvzID uuid.UUID) (*models.Reception, error) {
	return m.CloseLastReceptionFunc(pvzID)
}
func (m *mockPVZService) DeleteLastProduct(pvzID uuid.UUID) error {
	return m.DeleteLastProductFunc(pvzID)
}

func TestCreatePVZ(t *testing.T) {
	mockService := &mockPVZService{
		CreatePVZFunc: func(city string) (*models.PVZ, error) {
			return &models.PVZ{ID: uuid.New(), City: city}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	reqBody := `{"city": "Москва"}`
	r := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreatePVZ(w, r)
	res := w.Result()
	assert.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestCreatePVZ_BadRequest(t *testing.T) {
	mockService := &mockPVZService{
		CreatePVZFunc: func(city string) (*models.PVZ, error) {
			return &models.PVZ{ID: uuid.New(), City: city}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)

	reqBody := `{"city": "Москва"`
	r := httptest.NewRequest(http.MethodPost, "/pvz", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	h.CreatePVZ(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateReception(t *testing.T) {
	pvzID := uuid.New()
	mockService := &mockPVZService{
		CreateReceptionFunc: func(id uuid.UUID) (*models.Reception, error) {
			return &models.Reception{ID: uuid.New(), PVZID: id}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	body := map[string]string{"pvzId": pvzID.String()}
	buf, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	h.CreateReception(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateReception_BadRequest(t *testing.T) {
	mockService := &mockPVZService{
		CreateReceptionFunc: func(id uuid.UUID) (*models.Reception, error) {
			return &models.Reception{ID: uuid.New(), PVZID: id}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	body := map[string]string{"pvzId": ""}
	buf, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	h.CreateReception(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddProduct(t *testing.T) {
	pvzID := uuid.New()
	mockService := &mockPVZService{
		AddProductFunc: func(productType string, id uuid.UUID) (*models.Product, error) {
			return &models.Product{ID: uuid.New(), DateTime: time.Now(), Type: productType, ReceptionID: id}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	body := map[string]interface{}{"type": "book", "pvzId": pvzID.String()}
	buf, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	h.AddProduct(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteLastProduct(t *testing.T) {
	pvzID := uuid.New()
	mockService := &mockPVZService{
		DeleteLastProductFunc: func(id uuid.UUID) error {
			return nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	r := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": pvzID.String()})
	w := httptest.NewRecorder()

	h.DeleteLastProduct(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCloseLastReception(t *testing.T) {
	pvzID := uuid.New()
	mockService := &mockPVZService{
		CloseLastReceptionFunc: func(id uuid.UUID) (*models.Reception, error) {
			return &models.Reception{ID: uuid.New(), PVZID: id}, nil
		},
	}
	h := handlers.NewPVZHandler(mockService)
	r := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/close_last_reception", nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": pvzID.String()})
	w := httptest.NewRecorder()

	h.CloseLastReception(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetPVZList(t *testing.T) {
	mockService := &mockPVZService{
		GetPVZListFunc: func(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
			return []models.PVZWithReceptions{
				{
					PVZ: models.PVZ{
						ID:               uuid.New(),
						RegistrationDate: time.Now(),
						City:             "Москва",
					},
					Receptions: []models.ReceptionWithProducts{},
				},
			}, nil
		},
	}
	handler := handlers.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz/list?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	handler.GetPVZList(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var result []models.PVZWithReceptions
	err := json.NewDecoder(res.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Москва", result[0].PVZ.City)
}

func TestAddProduct_BadRequest(t *testing.T) {
	mockService := &mockPVZService{}
	h := handlers.NewPVZHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("{invalid-json"))
	w := httptest.NewRecorder()

	h.AddProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddProduct_InvalidUUID(t *testing.T) {
	mockService := &mockPVZService{}
	h := handlers.NewPVZHandler(mockService)

	body := map[string]interface{}{"type": "book", "pvzId": "not-a-uuid"}
	buf, _ := json.Marshal(body)
	r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	h.AddProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLastProduct_InvalidUUID(t *testing.T) {
	mockService := &mockPVZService{}
	h := handlers.NewPVZHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/delete_last_product", nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": "invalid-uuid"})
	w := httptest.NewRecorder()

	h.DeleteLastProduct(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReception_InvalidUUID(t *testing.T) {
	mockService := &mockPVZService{}
	h := handlers.NewPVZHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/close_last_reception", nil)
	r = mux.SetURLVars(r, map[string]string{"pvzId": "invalid-uuid"})
	w := httptest.NewRecorder()

	h.CloseLastReception(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCloseLastReception_NoUUID(t *testing.T) {
	mockService := &mockPVZService{}
	h := handlers.NewPVZHandler(mockService)

	r := httptest.NewRequest(http.MethodPost, "/pvz/invalid-uuid/close_last_reception", nil)
	r = mux.SetURLVars(r, map[string]string{"": ""})
	w := httptest.NewRecorder()

	h.CloseLastReception(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZList_InvalidPage(t *testing.T) {
	mockService := &mockPVZService{}
	handler := handlers.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz/list?page=abc", nil)
	w := httptest.NewRecorder()

	handler.GetPVZList(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZList_InvalidLimit(t *testing.T) {
	mockService := &mockPVZService{}
	handler := handlers.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz/list?limit=999", nil)
	w := httptest.NewRecorder()

	handler.GetPVZList(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZList_InvalidStartDate(t *testing.T) {
	mockService := &mockPVZService{}
	handler := handlers.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz/list?startDate=not-a-date", nil)
	w := httptest.NewRecorder()

	handler.GetPVZList(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPVZList_InvalidEndDate(t *testing.T) {
	mockService := &mockPVZService{}
	handler := handlers.NewPVZHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/pvz/list?endDate=not-a-date", nil)
	w := httptest.NewRecorder()

	handler.GetPVZList(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
