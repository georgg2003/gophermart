package postgres

import (
	"context"
	"database/sql"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/internal/usecase"
	"github.com/georgg2003/gophermart/pkg/errutils"
)

func (p *postgres) ApplyOrderAccrual(
	ctx context.Context,
	orderNumber string,
	accrual int,
) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection while applying an order accrual",
		)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return errutils.Wrap(err, "failed to begin a transaction while applying an order accrual")
	}
	defer tx.Rollback(ctx)

	var userID sql.NullInt64
	err = tx.QueryRow(
		ctx,
		`UPDATE orders SET status = $1 WHERE number = $2 RETURNING user_id`,
		models.StatusProcessed,
		orderNumber,
	).Scan(&userID)
	if err != nil {
		return errutils.Wrap(err, "failed to update order status")
	}
	if !userID.Valid {
		return usecase.ErrNoOrdersToUpdate
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO transactions (user_id, order_number, amount) VALUES ($1, $2, $3)`,
		userID.Int64,
		orderNumber,
		accrual,
	)
	if err != nil {
		return errutils.Wrap(err, "failed to insert transaction")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errutils.Wrap(err, "failed to commit transaction")
	}

	return nil
}
