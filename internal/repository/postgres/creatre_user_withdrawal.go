package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *postgres) CreateUserWithdrawal(
	ctx context.Context,
	userID int64,
	orderNumber string,
	amount int64,
) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when creating a new withdrawal",
		)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return errutils.Wrap(err, "failed to begin a transaction when creating a new withdrawal")
	}
	defer tx.Rollback(ctx)

	var balance int64
	err = tx.QueryRow(
		ctx,
		"SELECT coalesce(sum(amount), 0) FROM (SELECT * FROM transactions WHERE user_id = $1 FOR UPDATE)",
		userID,
	).Scan(&balance)

	if err != nil {
		return errutils.Wrap(err, "failed to select user balance")
	}
	if balance < amount {
		return usecase.ErrNotEnoughBalance
	}

	_, err = tx.Exec(
		ctx,
		"INSERT INTO transactions (user_id, order_number, amount) VALUES ($1, $2, $3)",
		userID, orderNumber, -amount,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return usecase.ErrWithdrawalAlreadyExists
		}
		return errutils.Wrap(err, "failed to insert withdrawal transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errutils.Wrap(err, "failed to commit tx")
	}

	return nil
}
