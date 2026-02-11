package core

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
