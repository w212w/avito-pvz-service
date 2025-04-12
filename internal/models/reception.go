package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `json:"id" db:"id"`
	PVZID    uuid.UUID `json:"pvz_id" db:"pvz_id"`
	DateTime time.Time `json:"date_time" db:"date_time"`
	Status   string    `json:"status" db:"status"`
}

type ReceptionWithProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}
