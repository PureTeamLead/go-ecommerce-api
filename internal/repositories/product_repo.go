package repositories

import (
	"database/sql"
	"errors"
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"eshop/pkg/postgre"
	"fmt"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *entities.Product) (uuid.UUID, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*entities.Product, error)
	GetAll() ([]entities.Product, error)
	Update(product *entities.Product) (*entities.Product, error)
}

// TODO: create model(request) to communicate with service

type productRepository struct {
	db postgre.DBinteraction
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (p *productRepository) Create(product *entities.Product) (uuid.UUID, error) {
	var id uuid.UUID
	newProduct := entities.NewProduct(product.Price, product.Name, product.Category)

	const query = `INSERT INTO products(id, price, category, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

	err := p.db.QueryRow(query, newProduct.ID, newProduct.Price, newProduct.Category, newProduct.Name, newProduct.CreatedAt, newProduct.UpdatedAt).Scan(&id)
	if err != nil {
		return constants.EmptyID, fmt.Errorf("failed creating new product: %w", err)
	}

	return id, nil
}

func (p *productRepository) Delete(id uuid.UUID) error {
	var returnedID uuid.UUID
	const query = `DELETE FROM products WHERE id = $1 RETURNING id;`

	err := p.db.QueryRow(query, id).Scan(&returnedID)
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNoProductFound
	}
	if (err != nil) || (returnedID != id) {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (p *productRepository) GetByID(id uuid.UUID) (*entities.Product, error) {
	var product entities.Product
	const query = `SELECT id, price, category, name, created_at, updated_at FROM products WHERE id = $1;`

	err := p.db.QueryRow(query, id).Scan(&product)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNoProductFound
	} else if err != nil {
		return nil, fmt.Errorf("refused to communicate with products database: %w", err)
	}

	return &product, nil
}

func (p *productRepository) GetAll() ([]entities.Product, error) {
	var products []entities.Product
	const query = `SELECT id, price, category, name, created_at, updated_at FROM products;`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product entities.Product
		err = rows.Scan(&product)
		if err != nil {
			return nil, fmt.Errorf("failed to scan structs: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate over products: %w", err)
	}

	return products, nil
}

func (p *productRepository) Update(product *entities.Product) (*entities.Product, error) {
	const query = `UPDATE products SET name = $1, price = $2, category = $3, updated_at = $4 WHERE id = $5;`
	updatedProduct := entities.UpdateProduct(product.ID, product.Name, product.Category, product.Price, product.CreatedAt)

	_, err := p.db.Exec(query, updatedProduct.Name, updatedProduct.Price, updatedProduct.Category, updatedProduct.UpdatedAt, updatedProduct.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNoProductFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedProduct, nil
}
