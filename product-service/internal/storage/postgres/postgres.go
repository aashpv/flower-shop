package postgres

import (
	"database/sql"
	"errors"
	"flower-shop/product-service/internal/model"
	"flower-shop/product-service/internal/storage"
	"fmt"
	_ "github.com/lib/pq" // driver for postgres
)

type Database struct {
	db *sql.DB
}

func New(storageConn string) (*Database, error) {
	db, err := sql.Open("postgres", storageConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Database{db: db}, nil
}

func (db *Database) GetAllProducts() ([]model.Product, error) {
	const op = "storage.postgres.GetAllProducts"
	var products []model.Product

	rows, err := db.db.Query(
		"SELECT name, description, price FROM products",
	)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, storage.ErrInternalServer)
	}
	for rows.Next() {
		var product model.Product

		if err := rows.Scan(
			&product.Name,
			&product.Description,
			&product.Price,
		); err != nil {
			return nil, fmt.Errorf("%s: scan failed: %w", op, storage.ErrInternalServer)
		}

		products = append(products, product)
	}

	return products, nil
}

// return id????
func (db *Database) CreateProduct(product *model.Product) (int, error) {
	const op = "storage.postgres.CreateProduct"
	var productId int

	err := db.db.QueryRow(
		"INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id",
		product.Name,
		product.Description,
		product.Price,
	).Scan(&productId)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create product: %w", op, storage.ErrInternalServer)
	}

	return productId, nil
}

// use?????
func (db *Database) GetProduct(id int) (*model.Product, error) {
	const op = "storage.postgres.GetProduct"
	var product model.Product

	err := db.db.QueryRow(
		"SELECT name, description, price FROM products WHERE id = $1",
		id,
	).Scan(
		&product.Name,
		&product.Description,
		&product.Price,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: sql error: %w", op, storage.ErrProductNotFound)
		}
		return nil, fmt.Errorf("%s: failed to get product: %w", op, storage.ErrInternalServer)
	}

	return &product, nil
}

func (db *Database) DeleteProduct(id int) error {
	const op = "storage.postgres.DeleteProduct"

	res, err := db.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: failed to delete product: %w", op, storage.ErrInternalServer)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, storage.ErrInternalServer)
	}
	if rows == 0 {
		return fmt.Errorf("%s: sql error: %w", op, storage.ErrProductNotFound)
	}

	return nil
}
