package services

import (
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PVZService interface {
	CreatePVZ(city string) (*models.PVZ, error)
	GetPVZList(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
	CreateReception(pvzID uuid.UUID) (*models.Reception, error)
	AddProduct(prodType string, pvzID uuid.UUID) (*models.Product, error)
	CloseLastReception(pvzID uuid.UUID) (*models.Reception, error)
	DeleteLastProduct(pvzID uuid.UUID) error
}

type pvzService struct {
	pvzRepo repository.PVZRepoInterface
}

func NewPVZService(pvzRepo repository.PVZRepoInterface) PVZService {
	return &pvzService{pvzRepo: pvzRepo}
}

func (s *pvzService) CreatePVZ(city string) (*models.PVZ, error) {

	validCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}
	if !validCities[city] {
		return &models.PVZ{}, fmt.Errorf("invalid city")
	}

	pvz := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             city,
	}

	err := s.pvzRepo.CreatePVZ(pvz)
	if err != nil {
		return nil, err
	}
	return pvz, nil
}

func (s *pvzService) GetPVZList(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	return s.pvzRepo.FetchPVZWithReceptions(startDate, endDate, page, limit)
}

var ErrReceptionExists = errors.New("reception already exists")

func (s *pvzService) CreateReception(pvzID uuid.UUID) (*models.Reception, error) {
	exists, err := s.pvzRepo.HasOpenReception(pvzID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrReceptionExists
	}

	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		DateTime: time.Now(),
		Status:   "in_progress",
	}

	if err := s.pvzRepo.CreateReception(reception); err != nil {
		return nil, err
	}

	return reception, nil
}

func (s *pvzService) AddProduct(prodType string, pvzID uuid.UUID) (*models.Product, error) {
	reception, err := s.pvzRepo.GetActiveReception(pvzID)
	if err != nil {
		return nil, errors.New("no active reception")
	}

	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now().UTC(),
		Type:        prodType,
		ReceptionID: reception.ID,
	}

	if err := s.pvzRepo.CreateProduct(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *pvzService) CloseLastReception(pvzID uuid.UUID) (*models.Reception, error) {
	reception, err := s.pvzRepo.GetLastReception(pvzID)
	if err != nil {
		return nil, err
	}

	if reception.Status == "close" {
		return nil, errors.New("reception already closed")
	}

	reception.Status = "close"

	err = s.pvzRepo.UpdateReceptionStatus(reception.ID, "close")
	if err != nil {
		return nil, err
	}
	return reception, nil

}

func (s *pvzService) DeleteLastProduct(pvzID uuid.UUID) error {
	reception, err := s.pvzRepo.GetLastReception(pvzID)
	if err != nil {
		return err
	}

	if reception.Status != "in_progress" {
		return errors.New("reception is not in progress")
	}

	product, err := s.pvzRepo.GetLastProduct(reception.ID)
	if err != nil {
		return err
	}

	return s.pvzRepo.DeleteProduct(product.ID)
}
