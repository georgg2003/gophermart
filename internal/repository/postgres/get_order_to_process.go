package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetOrderToProcess(
	ctx context.Context,
	processRetryTimeout int,
) (string, error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection while getting a new order",
		)
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", errutils.Wrap(err, "failed to begin a transaction while getting a new order")
	}
	defer tx.Rollback(ctx)

	var orderNumber string
	err = tx.QueryRow(
		ctx,
		`SELECT number FROM (
			SELECT number FROM orders 
			WHERE status = 'NEW' 
				OR (status = 'PROCESSING' AND processing_since < (NOW() - $1::interval))
			ORDER BY uploaded_at ASC
			LIMIT 1
		) FOR UPDATE SKIP LOCKED`,
		time.Duration(processRetryTimeout),
	).Scan(&orderNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", errutils.Wrap(err, "failed to select order to process")
	}

	if _, err = tx.Exec(
		ctx,
		`UPDATE orders
		SET status = 'PROCESSING', processing_since = NOW()
		WHERE number = $1`,
		orderNumber,
	); err != nil {
		return "", errutils.Wrap(err, "failed to update order to process")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", errutils.Wrap(err, "failed to commit tx")
	}
	return orderNumber, nil
}
