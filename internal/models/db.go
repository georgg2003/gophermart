package models

import "time"

type WithdrawalDB struct {
	Order       *string
	Amount      int64
	ProcessedAt *time.Time `db:"processed_at"`
}

type OrderDB struct {
	Number     string
	Status     OrderStatus
	Accrual    int64
	UploadedAt time.Time `db:"uploaded_at"`
}
