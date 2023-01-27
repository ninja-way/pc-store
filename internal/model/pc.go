package model

type PC struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CPU         string `json:"cpu"`
	Videocard   string `json:"videocard"`
	RAM         int    `json:"ram"`
	DataStorage string `json:"data_storage"`
	Price       int    `json:"price"`
}
