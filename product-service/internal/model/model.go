package model

type Product struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Price       int    `json:"price" validate:"required,numeric,min=0"`
	Description string `json:"description" validate:"required,max=500"`
}
