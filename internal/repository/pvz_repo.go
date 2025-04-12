package repository

import (
	"avito-pvz-service/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PVZRepoInterface interface {
	CreatePVZ(pvz *models.PVZ) error
	FetchPVZWithReceptions(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
}

type PVZRepo struct {
	db *sqlx.DB
}

func NewPVZRepo(db *sql.DB) PVZRepoInterface {
	return &PVZRepo{db: sqlx.NewDb(db, "postgres")}
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

func (r *PVZRepo) FetchPVZWithReceptions(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	offset := (page - 1) * limit

	query := `SELECT * FROM pvz LIMIT $1 OFFSET $2`
	var pvzs []models.PVZ
	if err := r.db.Select(&pvzs, query, limit, offset); err != nil {
		return nil, errors.New("error fetch")
	}

	result := make([]models.PVZWithReceptions, 0)

	for _, pvz := range pvzs {
		receptions := []models.ReceptionWithProducts{}
		var rList []models.Reception

		rQuery := `SELECT * FROM receptions WHERE pvz_id = $1`
		args := []interface{}{pvz.ID}
		argIdx := 2 // $1 — pvz.ID

		if startDate != nil {
			rQuery += fmt.Sprintf(" AND date_time >= $%d", argIdx)
			args = append(args, *startDate)
			argIdx++
		}
		if endDate != nil {
			rQuery += fmt.Sprintf(" AND date_time <= $%d", argIdx)
			args = append(args, *endDate)
			argIdx++
		}

		if err := r.db.Select(&rList, rQuery, args...); err != nil {
			return nil, errors.New("error fetch")
		}

		for _, reception := range rList {
			var products []models.Product
			rpQuery := `SELECT * FROM products WHERE reception_id = $1`
			if err := r.db.Select(&products, rpQuery, reception.ID); err != nil {
				return nil, errors.New("error fetch")
			}

			receptions = append(receptions, models.ReceptionWithProducts{
				Reception: reception,
				Products:  products,
			})
		}

		result = append(result, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: receptions,
		})
	}

	return result, nil
}
