package models

import "time"

func NewWithdrawalFromDB(db WithdrawalDB) Withdrawal {
	var order *string
	if db.Order.Valid {
		order = &db.Order.String
	} else {
		order = nil
	}

	var processedAt *time.Time
	if db.ProcessedAt.Valid {
		processedAt = &db.ProcessedAt.Time
	} else {
		processedAt = nil
	}

	return Withdrawal{
		Order:       order,
		Sum:         NewMoney(db.Amount),
		ProcessedAt: processedAt,
	}
}

func NewWithdrawalsFromDB(db []WithdrawalDB) []Withdrawal {
	res := make([]Withdrawal, 0, len(db))
	for _, v := range db {
		res = append(res, NewWithdrawalFromDB(v))
	}
	return res
}

func NewOrderFromDB(db OrderDB) Order {
	var accrual *Money
	if db.Accrual.Valid {
		money := NewMoney(db.Accrual.Int64)
		accrual = &money
	} else {
		accrual = nil
	}
	return Order{
		Number:     db.Number,
		Status:     db.Status,
		Accrual:    accrual,
		UploadedAt: db.UploadedAt,
	}
}

func NewOrdersFromDB(db []OrderDB) []Order {
	res := make([]Order, 0, len(db))
	for _, v := range db {
		res = append(res, NewOrderFromDB(v))
	}
	return res
}
