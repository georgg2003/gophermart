package models

type UserCredentials struct {
	ID           int64
	Login        string
	PasswordHash string
}
