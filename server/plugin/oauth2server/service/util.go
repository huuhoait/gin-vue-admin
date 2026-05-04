package service

import (
	"crypto/rand"
	"encoding/hex"
)

// secureToken returns a hex-encoded random byte string. 32 bytes → 64-char
// token; sufficient entropy for opaque OAuth2 access/refresh tokens.
func secureToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
