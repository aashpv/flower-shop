package handlers

import "flower-shop/product-service/internal/model"

type Handlers interface {
	CreateProduct(product *model.Product) (int, error)
	GetAllProducts() ([]model.Product, error)
	GetProduct(id int) (*model.Product, error)
	DeleteProduct(id int) error
}
