package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id" db:"id"`
	DateTime    time.Time `json:"date_time" db:"date_time"`
	Type        string    `json:"type" db:"type"`
	ReceptionID uuid.UUID `json:"reception_id" db:"reception_id"`
}
