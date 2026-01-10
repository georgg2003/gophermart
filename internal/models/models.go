package models

import (
	"time"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
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

type Order struct {
	Number     string
	Status     OrderStatus
	Accrual    *Money
	UploadedAt time.Time
}
