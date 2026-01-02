package models

import (
	"time"
)

type UserCredentials struct {
	ID           int64
	Login        string
	PasswordHash string
}

type UserBalance struct {
	Current   Money
	Withdrawn Money
}

type Withdrawal struct {
	Order       *string
	Sum         Money
	ProcessedAt *time.Time
}

type WithdrawalDB struct {
	Order       *string
	Amount      int64
	ProcessedAt *time.Time `db:"processed_at"`
}

func NewWithdrawalFromDB(db WithdrawalDB) Withdrawal {
	return Withdrawal{
		Order:       db.Order,
		Sum:         NewMoney(db.Amount),
		ProcessedAt: db.ProcessedAt,
	}
}

func NewWithdrawalsFromDB(db []WithdrawalDB) []Withdrawal {
	res := make([]Withdrawal, 0, len(db))
	for _, v := range db {
		res = append(res, NewWithdrawalFromDB(v))
	}
	return res
}
