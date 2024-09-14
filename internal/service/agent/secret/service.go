package secret

import (
	"crypto/sha256"
	"fmt"
)

func GetHash(key []byte) {

	fmt.Println(key)

	// создаём новый hash.Hash, вычисляющий контрольную сумму SHA-256
	h := sha256.New()
	// передаём байты для хеширования
	h.Write(key)
	// вычисляем хеш
	dst := h.Sum(nil)

	fmt.Printf("%x", dst)
}
