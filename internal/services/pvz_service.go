package services

import (
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PVZService struct {
	pvzRepo repository.PVZRepoInterface
}

func NewPVZService(pvzRepo repository.PVZRepoInterface) *PVZService {
	return &PVZService{pvzRepo: pvzRepo}
}

func (s *PVZService) CreatePVZ(city string) (*models.PVZ, error) {

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

func (s *PVZService) GetPVZList(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	return s.pvzRepo.FetchPVZWithReceptions(startDate, endDate, page, limit)
}
