package handlers

import "flower-shop/user-service/internal/model"

type Handlers interface {
	CreateUser(user *model.User) (int, error)
	CheckUser(email string) (int, string, string, error)
	DeleteUser(email string) error
}
