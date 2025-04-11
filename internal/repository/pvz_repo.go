package repository

import (
	"avito-pvz-service/internal/models"
	"database/sql"
	"errors"
)

type PVZRepoInterface interface {
	CreatePVZ(pvz *models.PVZ) error
}

type PVZRepo struct {
	db *sql.DB
}

func NewPVZRepo(db *sql.DB) PVZRepoInterface {
	return &PVZRepo{db: db}
}

func (r *PVZRepo) CreatePVZ(pvz *models.PVZ) error {
	query := `
	INSERT INTO pvz (id, registration_date, city)
	VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query,
		pvz.ID,
		pvz.RegistrationDate,
		pvz.City,
	)

	if err != nil {
		return errors.New("error to create pvz")
	}
	return nil
}
