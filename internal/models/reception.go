package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	PVZID    uuid.UUID `json:"pvz_id"`
	DateTime time.Time `json:"date_time"`
	Status   string    `json:"status"`
}
