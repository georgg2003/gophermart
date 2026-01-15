package usecase

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserWrongPassword = errors.New("password is incorrect")
var ErrUserNotFound = errors.New("user not found")

var ErrWidthdrawalsNotFound = errors.New("withdrawals not found")

var ErrOrdersNotFound = errors.New("orders not found")

var ErrOrderAlreadyUploaded = errors.New("order already uploaded")
var ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")

var ErrNotEnoughBalance = errors.New("not enough balance to withdraw")
var ErrWithdrawalAlreadyExists = errors.New("withdrawal already exists")

var ErrNoOrdersToUpdate = errors.New("no orders to update")
