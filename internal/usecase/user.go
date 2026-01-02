package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/georgg2003/gophermart/internal/models"
	"github.com/georgg2003/gophermart/pkg/errutils"
)

func passwordToHash(password string) string {
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(password)))
}

func (uc *useCase) UserRegister(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := passwordToHash(password)
	userID, err := uc.repo.NewUser(ctx, login, hash)
	if err != nil {
		return "", errutils.Wrap(err, "failed to create new user")
	}

	return uc.jwtHelper.NewAccessToken(userID)
}

func (uc *useCase) UserLogin(
	ctx context.Context,
	login string,
	password string,
) (accessToken string, err error) {
	hash := passwordToHash(password)
	userCredentials, err := uc.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", errutils.Wrap(err, "failed to get user by login")
	}
	if userCredentials.PasswordHash != hash {
		return "", ErrUserWrongPassword
	}

	return uc.jwtHelper.NewAccessToken(userCredentials.ID)
}

func (uc *useCase) UserGetBalance(
	ctx context.Context,
	userID int64,
) (balance *models.UserBalance, err error) {
	balance, err = uc.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, errutils.Wrap(err, "failed to get user balance")
	}
	return balance, err
}
