package model

import "time"

type PC struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CPU         string    `json:"cpu"`
	Videocard   string    `json:"videocard,omitempty"`
	RAM         int       `json:"ram"`
	DataStorage string    `json:"data_storage,omitempty"`
	AddedAt     time.Time `json:"added_at"`
	Price       int       `json:"price"`
}
