package services

import (
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"github.com/google/uuid"
)

type productRepository interface {
	Create(product *entities.Product) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entities.Product, error)
	GetAll() ([]entities.Product, error)
	Update(product *entities.Product) error
}

type ProductService struct {
	repo productRepository
}

func NewProductService(r productRepository) *ProductService {
	return &ProductService{repo: r}
}

func (ps *ProductService) AddProduct(name, category string, price float64) (uuid.UUID, error) {
	newProduct := entities.NewProduct(price, name, category)
	id, err := ps.repo.Create(newProduct)
	if err != nil {
		return constants.EmptyID, err
	}

	return id, nil
}

func (ps *ProductService) DeleteProduct(id uuid.UUID) (uuid.UUID, error) {
	if err := ps.repo.Delete(id); err != nil {
		return constants.EmptyID, err
	}

	return id, nil
}

func (ps *ProductService) UpdateProduct(id uuid.UUID, name, category string, price float64) (*entities.Product, error) {
	product, err := ps.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedProduct := entities.UpdateProduct(id, name, category, price, product.CreatedAt)

	if err = ps.repo.Update(updatedProduct); err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *ProductService) GetProduct(id uuid.UUID) (*entities.Product, error) {
	product, err := ps.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductService) GetAllProducts() ([]entities.Product, error) {
	products, err := ps.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}
