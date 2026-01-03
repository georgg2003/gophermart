package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *postgres) CreateUserOrder(
	ctx context.Context,
	userID int64,
	orderNumber string,
) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when creating a new user order",
		)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		"INSERT INTO orders (user_id, number) VALUES ($1, $2)",
		userID, orderNumber,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			var isOwnOrder bool
			conn.QueryRow(
				ctx,
				"SELECT EXISTS (SELECT 1 FROM orders WHERE number = $1 AND user_id = $2)",
				orderNumber,
				userID,
			).Scan(&isOwnOrder)
			if isOwnOrder {
				return usecase.ErrOrderAlreadyUploaded
			} else {
				return usecase.ErrOrderAlreadyUploadedByAnotherUser
			}
		}
		return errutils.Wrap(err, "failed to insert a new user order")
	}

	return nil
}
