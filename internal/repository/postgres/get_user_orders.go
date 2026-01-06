package postgres

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetUserOrders(
	ctx context.Context,
	userID int64,
) ([]models.Order, error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when getting user orders",
		)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			ord.number as number,
			ord.status as status,
			tr.amount as accrual,
			ord.uploaded_at as uploaded_at
		FROM orders ord
		LEFT JOIN transactions tr on ord.number = tr.order_number
		WHERE ord.user_id = $1
		ORDER BY ord.uploaded_at DESC`,
		userID,
	)
	if err != nil {
		err = errutils.Wrap(err, "failed to select user orders")
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.OrderDB])

	return models.NewOrdersFromDB(orders), err
}
