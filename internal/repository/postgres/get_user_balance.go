package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/pkg/errutils"
	"github.com/jackc/pgx/v5"
)

func (p *postgres) GetUserBalance(
	ctx context.Context,
	userID int64,
) (*models.UserBalance, error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		err = errutils.Wrap(
			err,
			"failed to acquire a db connection when getting a user by login",
		)
		return nil, err
	}
	defer conn.Release()

	var balance models.UserBalance
	err = conn.QueryRow(
		ctx,
		`SELECT 
			COALESCE(SUM(amount), 0) AS current,
			COALESCE(SUM(
				CASE
        	WHEN amount < 0 THEN amount
    		END
			), 0) AS withdrawn 
			FROM transactions 
			WHERE user_id = $1`,
		userID,
	).Scan(&balance.Current.AmountMinor, &balance.Withdrawn.AmountMinor)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &balance, nil
		}
		err = errutils.Wrap(err, "failed to get user balance")
		return nil, err
	}

	return &balance, err
}
