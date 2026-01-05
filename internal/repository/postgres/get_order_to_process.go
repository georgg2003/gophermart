package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/georgg2003/gophermart/pkg/errutils"
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

	var orderNumber sql.NullString
	err = tx.QueryRow(
		ctx,
		`SELECT number FROM (
			SELECT number FROM orders 
			WHERE status = 'NEW' 
				OR (status = 'PROCESSING' AND NOW() - processing_since > INTERVAL $1) 
			LIMIT 1
		) FOR UPDATE SKIP LOCKED`,
		time.Duration(processRetryTimeout),
	).Scan(&orderNumber)
	if err != nil {
		return "", errutils.Wrap(err, "failed to get order to process")
	}
	if !orderNumber.Valid {
		return "", nil
	}
	orderNumberStr := orderNumber.String

	if _, err = tx.Exec(
		ctx,
		`UPDATE orders
		SET status = 'PROCESSING', processing_since = NOW()
		WHERE number = $1`,
		orderNumberStr,
	); err != nil {
		return "", errutils.Wrap(err, "failed to update order to process")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", errutils.Wrap(err, "failed to commit tx")
	}
	return orderNumberStr, nil
}
