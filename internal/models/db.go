package models

import (
	"database/sql"
	"time"
)

type WithdrawalDB struct {
	Order       sql.NullString
	Amount      int64
	ProcessedAt sql.NullTime `db:"processed_at"`
}

type OrderDB struct {
	Number     string
	Status     OrderStatus
	Accrual    sql.NullInt64
	UploadedAt time.Time `db:"uploaded_at"`
}
