package models

import "time"

type Reception struct {
	ID       string    `json:"id"`
	PVZID    string    `json:"pvz_id"`
	DateTime time.Time `json:"date_time"`
	Status   string    `json:"status"`
}
