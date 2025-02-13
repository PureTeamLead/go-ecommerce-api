package services

import "eshop/internal/repositories"

type ProductService interface {
	AddProduct()
	DeleteProduct()
	UpdateProduct()
	GetProduct()
	GetAllProducts()
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(r repositories.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (ps *productService) AddProduct() {

}

func (ps *productService) DeleteProduct() {

}

func (ps *productService) UpdateProduct() {

}

func (ps *productService) GetProduct() {

}

func (ps *productService) GetAllProducts() {

}
