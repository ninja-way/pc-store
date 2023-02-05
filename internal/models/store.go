package models

type Store struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Computers []PC   `json:"computers"`
}
