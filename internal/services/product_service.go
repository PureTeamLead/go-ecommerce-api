package services

import (
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/repositories"
	"github.com/google/uuid"
)

type ProductService interface {
	AddProduct(product *entities.Product) (uuid.UUID, error)
	DeleteProduct(id uuid.UUID) (uuid.UUID, error)
	UpdateProduct(product *entities.Product) (*entities.Product, error)
	GetProduct(id uuid.UUID) (*entities.Product, error)
	GetAllProducts() ([]entities.Product, error)
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(r repositories.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (ps *productService) AddProduct(product *entities.Product) (uuid.UUID, error) {
	id, err := ps.repo.Create(product)
	if err != nil {
		return constants.EmptyID, err
	}

	return id, nil
}

func (ps *productService) DeleteProduct(id uuid.UUID) (uuid.UUID, error) {
	if err := ps.repo.Delete(id); err != nil {
		return constants.EmptyID, err
	}

	return id, nil
}

func (ps *productService) UpdateProduct(product *entities.Product) (*entities.Product, error) {
	updatedProduct := entities.UpdateProduct(product.ID, product.Name, product.Category, product.Price, product.CreatedAt)

	if err := ps.repo.Update(updatedProduct); err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (ps *productService) GetProduct(id uuid.UUID) (*entities.Product, error) {
	product, err := ps.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *productService) GetAllProducts() ([]entities.Product, error) {
	products, err := ps.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}
