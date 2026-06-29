package refreshtoken

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// generate token
func Generate() string {
	return rand.Text()
}

// hashing token
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
