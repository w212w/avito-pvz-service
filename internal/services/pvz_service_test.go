package services_test

import (
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPVZRepo struct {
	mock.Mock
}

func (m *mockPVZRepo) CreatePVZ(pvz *models.PVZ) error {
	args := m.Called(pvz)
	return args.Error(0)
}

func (m *mockPVZRepo) FetchPVZWithReceptions(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	args := m.Called(startDate, endDate, page, limit)
	return args.Get(0).([]models.PVZWithReceptions), args.Error(1)
}

func (m *mockPVZRepo) HasOpenReception(pvzID uuid.UUID) (bool, error) {
	args := m.Called(pvzID)
	return args.Bool(0), args.Error(1)
}

func (m *mockPVZRepo) CreateReception(r *models.Reception) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *mockPVZRepo) GetActiveReception(pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *mockPVZRepo) CreateProduct(p *models.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *mockPVZRepo) GetLastReception(pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *mockPVZRepo) UpdateReceptionStatus(receptionID uuid.UUID, status string) error {
	args := m.Called(receptionID, status)
	return args.Error(0)
}

func (m *mockPVZRepo) GetLastProduct(receptionID uuid.UUID) (*models.Product, error) {
	args := m.Called(receptionID)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *mockPVZRepo) DeleteProduct(productID uuid.UUID) error {
	args := m.Called(productID)
	return args.Error(0)
}

func TestCreatePVZ_ValidCity(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)

	repo.On("CreatePVZ", mock.Anything).Return(nil)

	pvz, err := service.CreatePVZ("Москва")
	assert.NoError(t, err)
	assert.Equal(t, "Москва", pvz.City)
}

func TestCreatePVZ_InvalidCity(t *testing.T) {
	service := services.NewPVZService(new(mockPVZRepo))

	_, err := service.CreatePVZ("Новосибирск")
	assert.Error(t, err)
	assert.Equal(t, "invalid city", err.Error())
}

func TestCreateReception_WhenNoneExists(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)
	pvzID := uuid.New()

	repo.On("HasOpenReception", pvzID).Return(false, nil)
	repo.On("CreateReception", mock.Anything).Return(nil)

	r, err := service.CreateReception(pvzID)
	assert.NoError(t, err)
	assert.Equal(t, pvzID, r.PVZID)
	assert.Equal(t, "in_progress", r.Status)
}

func TestCreateReception_WhenAlreadyExists(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)
	pvzID := uuid.New()

	repo.On("HasOpenReception", pvzID).Return(true, nil)

	_, err := service.CreateReception(pvzID)
	assert.Error(t, err)
	assert.Equal(t, services.ErrReceptionExists, err)
}

func TestAddProduct_WithActiveReception(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)
	pvzID := uuid.New()
	reception := &models.Reception{ID: uuid.New()}

	repo.On("GetActiveReception", pvzID).Return(reception, nil)
	repo.On("CreateProduct", mock.Anything).Return(nil)

	product, err := service.AddProduct("electronics", pvzID)
	assert.NoError(t, err)
	assert.Equal(t, reception.ID, product.ReceptionID)
	assert.Equal(t, "electronics", product.Type)
}

func TestCloseLastReception(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)
	pvzID := uuid.New()
	reception := &models.Reception{ID: uuid.New(), Status: "in_progress"}

	repo.On("GetLastReception", pvzID).Return(reception, nil)
	repo.On("UpdateReceptionStatus", reception.ID, "close").Return(nil)

	r, err := service.CloseLastReception(pvzID)
	assert.NoError(t, err)
	assert.Equal(t, "close", r.Status)
}

func TestDeleteLastProduct(t *testing.T) {
	repo := new(mockPVZRepo)
	service := services.NewPVZService(repo)
	pvzID := uuid.New()
	reception := &models.Reception{ID: uuid.New(), Status: "in_progress"}
	product := &models.Product{ID: uuid.New()}

	repo.On("GetLastReception", pvzID).Return(reception, nil)
	repo.On("GetLastProduct", reception.ID).Return(product, nil)
	repo.On("DeleteProduct", product.ID).Return(nil)

	err := service.DeleteLastProduct(pvzID)
	assert.NoError(t, err)
}
