package postgres

import (
	"database/sql"
	"errors"
	"flower-shop/user-service/internal/model"
	"flower-shop/user-service/internal/storage"
	"fmt"
	"github.com/lib/pq" // driver for postgres
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

func (db *Database) CreateUser(user *model.User) (int, error) {
	const op = "storage.postgres.CreateUser"
	var userId int

	err := db.db.QueryRow(
		"INSERT INTO users (first_name, last_name, phone, address, role, email, password) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Address,
		user.Role,
		user.Email,
		user.Password,
	).Scan(&userId)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return 0, fmt.Errorf("%s: failed to create user: %w", op, storage.ErrUserAlreadyExists)
		}

		return 0, fmt.Errorf("%s: failed to create user: %w", op, err)
	}

	return userId, nil
}

func (db *Database) CheckUser(email string) (int, string, string, error) {
	const op = "storage.postgres.CheckUser"

	var id int
	var password string
	var role string

	err := db.db.QueryRow(
		"SELECT id, password, role FROM users WHERE email = $1",
		email,
	).Scan(&id, &password, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", "", storage.ErrUserNotFound
		}
		return 0, "", "", fmt.Errorf("%s: failed to check user: %w", op, err)
	}

	return id, password, role, nil
}

// dont use
func (db *Database) DeleteUser(email string) error {
	const op = "storage.postgres.DeleteUser"

	res, err := db.db.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		return fmt.Errorf("%s: execute err: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected err: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return nil
}

// dont use
func (db *Database) UpdateRole(email, role string) error {
	const op = "storage.postgres.UpdateRole"

	res, err := db.db.Exec("UPDATE users SET role = $1 WHERE email = $2", role, email)
	if err != nil {
		return fmt.Errorf("%s: execute err: %w", op, err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: rows affected err: %w", op, err)
	}
	if rows == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return nil
}
