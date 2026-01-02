package postgres

import (
	"context"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetUserWithdrawals(
	ctx context.Context,
	userID int64,
) ([]models.Withdrawal, error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when getting user withdrawals",
		)
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			ord.number as order,
			tr.amount as amount,
			ord.processed_at as processed_at
		FROM transactions tr
		LEFT JOIN orders ord on ord.id = tr.order_id
		WHERE tr.user_id = $1
		ORDER BY ord.uploaded_at DESC`,
		userID,
	)
	if err != nil {
		err = errutils.Wrap(err, "failed to get user withdrawals")
		return nil, err
	}

	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.WithdrawalDB])

	return models.NewWithdrawalsFromDB(withdrawals), err
}
