package repository

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
)

//go:generate go tool mockgen -destination ./mock/mock.go -package mock . Repository
type Repository interface {
	NewUser(
		ctx context.Context,
		login string,
		passwordHash string,
	) (id int64, err error)
	GetUserByLogin(
		ctx context.Context,
		login string,
	) (creds *models.UserCredentials, err error)
	GetUserBalance(
		ctx context.Context,
		userID int64,
	) (balance *models.UserBalance, err error)
}
