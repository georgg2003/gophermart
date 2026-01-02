package models

type UserCredentials struct {
	ID           int64
	Login        string
	PasswordHash string
}

type UserBalance struct {
	Current   Money
	Withdrawn Money
}
