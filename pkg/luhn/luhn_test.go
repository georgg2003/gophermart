package luhn_test

import (
	"testing"

	"github.com/georgg2003/gophermart/pkg/luhn"
	"github.com/stretchr/testify/assert"
)

func TestValidLuhn(t *testing.T) {
	assert.True(t, luhn.ValidLuhn("12345678903"))
	assert.False(t, luhn.ValidLuhn("123456789"))
}
