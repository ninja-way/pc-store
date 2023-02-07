package models

import "time"

type PC struct {
	ID          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"Super PC"`
	CPU         string    `json:"cpu" example:"i9"`
	Videocard   string    `json:"videocard,omitempty" example:"RTX"`
	RAM         int       `json:"ram" example:"32"`
	DataStorage string    `json:"data_storage,omitempty" example:"ssd 1tb"`
	AddedAt     time.Time `json:"added_at" example:"2023-01-01T00:00:00.00000Z"`
	Price       int       `json:"price" example:"79999"`
}
