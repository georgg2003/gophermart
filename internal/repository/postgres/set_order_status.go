package postgres

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/pkg/errutils"
)

func (p *postgres) SetOrderStatus(
	ctx context.Context,
	orderNumber string,
	orderStatus models.OrderStatus,
) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection while setting order status",
		)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`UPDATE orders SET status = $1 WHERE number = $2`,
		orderStatus,
		orderNumber,
	)
	if err != nil {
		return errutils.Wrap(err, "failed to update order status")
	}

	return nil
}
