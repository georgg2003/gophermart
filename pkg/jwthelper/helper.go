package jwthelper

import (
	"github.com/golang-jwt/jwt/v4"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

type JWTHelper struct {
	key []byte
}

func (h *JWTHelper) NewAccessToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		TokenClaims{
			UserID: userID,
		},
	)
	return token.SignedString(h.key)
}

func (h *JWTHelper) ReadAccessToken(encodedToken string) (int64, error) {
	parser := jwt.NewParser()

	claims := TokenClaims{}
	token, err := parser.ParseWithClaims(encodedToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return h.key, jwt.NewValidationError(
				"signing method is not correct",
				jwt.ValidationErrorUnverifiable,
			)
		}

		return h.key, nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, jwt.ErrTokenNotValidYet
	}

	return claims.UserID, nil
}

func New(key []byte) *JWTHelper {
	return &JWTHelper{
		key: key,
	}
}
