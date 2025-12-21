package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *postgres) NewUser(
	ctx context.Context,
	login string,
	passwordHash string,
) (id int64, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when creating a new user",
		)
		return -1, err
	}
	defer conn.Release()

	var userID int64

	err = conn.QueryRow(
		ctx,
		"INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id",
		login, passwordHash,
	).Scan(&userID)

	if err == nil {
		return userID, nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return -1, usecase.ErrUserAlreadyExists
	}

	return -1, errutils.Wrap(err, "failed to insert a new user")
}
