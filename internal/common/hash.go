package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256sum(value []byte, key string) string {
	h := sha256.New()
	h.Write(value)
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
