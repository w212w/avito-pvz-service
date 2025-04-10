package models

import "time"

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	City             string    `json:"city"`
}
