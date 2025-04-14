package integration_test

import (
	"avito-pvz-service/config"
	"avito-pvz-service/internal/models"
	"avito-pvz-service/internal/repository"
	"avito-pvz-service/internal/storage"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var pvzRepo repository.PVZRepoInterface

func TestMain(m *testing.M) {
	cfg := config.LoadConfig()
	testdb := storage.ConnectDB(cfg)
	defer testdb.Close()
	pvzRepo = repository.NewPVZRepo(testdb)
	os.Exit(m.Run())
}

func TestIntegration_PVZReceptionFlow(t *testing.T) {
	pvz := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             "Москва",
	}
	err := pvzRepo.CreatePVZ(pvz)
	assert.NoError(t, err)

	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    pvz.ID,
		DateTime: time.Now(),
		Status:   "in_progress",
	}
	err = pvzRepo.CreateReception(reception)
	assert.NoError(t, err)

	for i := 0; i < 50; i++ {
		product := &models.Product{
			ID:          uuid.New(),
			DateTime:    time.Now(),
			Type:        "одежда",
			ReceptionID: reception.ID,
		}
		err := pvzRepo.CreateProduct(product)
		assert.NoError(t, err)
	}

	err = pvzRepo.UpdateReceptionStatus(reception.ID, "close")
	assert.NoError(t, err)

	updatedReception, err := pvzRepo.GetLastReception(pvz.ID)
	assert.NoError(t, err)
	assert.Equal(t, "close", updatedReception.Status)
}
