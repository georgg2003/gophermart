package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetUserByLogin(
	ctx context.Context,
	login string,
) (creds *models.UserCredentials, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when getting a user by login",
		)
		return nil, err
	}
	defer conn.Release()

	err = conn.QueryRow(
		ctx,
		"SELECT id, login, password_hash FROM users WHERE login = $1",
		login,
	).Scan(&creds.ID, &creds.Login, &creds.PasswordHash)

	if err == nil {
		return creds, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, usecase.ErrUserNotFound
	}

	err = errutils.Wrap(err, "failed to get a user by login")
	return nil, err
}
