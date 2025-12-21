package jwthelper

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

var testUserID int64 = 1
var key = []byte("secret_token")

func TestNewAccessToken(t *testing.T) {
	helper := New(key)
	token, err := helper.NewAccessToken(testUserID)
	assert.NoError(t, err)

	gotUserID, err := helper.ReadAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, testUserID, gotUserID)
}

func TestInvalidSigningMethod(t *testing.T) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		TokenClaims{
			UserID: testUserID,
		},
	)
	encodedToken, err := token.SignedString(key)
	assert.NoError(t, err)

	helper := New(key)
	_, err = helper.ReadAccessToken(encodedToken)

	var tgt *jwt.ValidationError
	assert.ErrorAs(t, err, &tgt)
	assert.True(t, tgt.Is(jwt.ErrTokenUnverifiable))
}
