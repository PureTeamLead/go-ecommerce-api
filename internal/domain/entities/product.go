package entities

import (
	"github.com/google/uuid"
	"time"
)

//Getters, setters, validating, constructor
// Sometimes need to create validated_product and constructor for it

type Product struct {
	ID        uuid.UUID
	Price     float64
	Name      string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(price float64, name, category string) *Product {
	return &Product{
		ID:        uuid.New(),
		Price:     price,
		Name:      name,
		Category:  category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func UpdateProduct(id uuid.UUID, name, category string, price float64, createdAt time.Time) *Product {
	return &Product{
		ID:        id,
		Name:      name,
		Category:  category,
		Price:     price,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}
}
