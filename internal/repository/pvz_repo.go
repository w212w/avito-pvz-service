package repository

import (
	"avito-pvz-service/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PVZRepoInterface interface {
	CreatePVZ(pvz *models.PVZ) error
	FetchPVZWithReceptions(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
	CreateReception(rec *models.Reception) error
	HasOpenReception(pvzID uuid.UUID) (bool, error)
	CreateProduct(p *models.Product) error
	GetActiveReception(pvzID uuid.UUID) (*models.Reception, error)
	GetLastReception(pvzID uuid.UUID) (*models.Reception, error)
	UpdateReceptionStatus(receptionID uuid.UUID, status string) error
	GetLastProduct(receptionID uuid.UUID) (*models.Product, error)
	DeleteProduct(productID uuid.UUID) error
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

func (r *PVZRepo) CreateReception(rec *models.Reception) error {
	query := `
		INSERT INTO receptions (id, pvz_id, date_time, status)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, rec.ID, rec.PVZID, rec.DateTime, rec.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *PVZRepo) HasOpenReception(pvzID uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
	`
	err := r.db.Get(&count, query, pvzID)
	if err != nil {
		return false, fmt.Errorf("failed to check open receptions: %w", err)
	}
	return count > 0, nil
}

func (r *PVZRepo) CreateProduct(p *models.Product) error {
	query := `INSERT INTO products (id, date_time, type, reception_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, p.ID, p.DateTime, p.Type, p.ReceptionID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PVZRepo) GetActiveReception(pvzID uuid.UUID) (*models.Reception, error) {
	var reception models.Reception
	query := `SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' LIMIT 1`
	err := r.db.Get(&reception, query, pvzID)
	if err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *PVZRepo) GetLastReception(pvzID uuid.UUID) (*models.Reception, error) {
	var reception models.Reception
	query := `SELECT id, date_time, pvz_id, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' LIMIT 1`
	err := r.db.Get(&reception, query, pvzID)
	if err != nil {
		return nil, err
	}
	return &reception, nil
}

func (r *PVZRepo) UpdateReceptionStatus(receptionID uuid.UUID, status string) error {
	query := `UPDATE receptions SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, receptionID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PVZRepo) GetLastProduct(receptionID uuid.UUID) (*models.Product, error) {
	product := models.Product{}
	query := `SELECT * FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1`
	err := r.db.Get(&product, query, receptionID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *PVZRepo) DeleteProduct(productID uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, productID)
	if err != nil {
		return err
	}
	return nil
}
