package requests

import "github.com/google/uuid"

type UpdateProduct struct {
	ID       uuid.UUID `json:"id"`
	Price    float64   `json:"price"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
}

type AddProduct struct {
	Price    float64 `json:"price"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
}

type DeleteGetProduct struct {
	ID uuid.UUID `json:"id"`
}
