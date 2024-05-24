package storage

import "errors"

var (
	// TODO: add errors
	ErrInternalServer  = errors.New("internal server error")
	ErrProductNotFound = errors.New("product not found")
)
