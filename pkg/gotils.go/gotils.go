package gotils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := sha256.New()
	return hex.EncodeToString(hash.Sum([]byte(password)))
}
