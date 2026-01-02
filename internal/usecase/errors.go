package usecase

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserWrongPassword = errors.New("password is incorrect")
var ErrUserNotFound = errors.New("user not found")

var ErrWidthdrawalsNotFound = errors.New("withdrawals not found")
