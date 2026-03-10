package utils

import (
	"crypto/sha256"
	"math/big"
)

const (
	maxLength        = 10
	availableSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base             = int64(len(availableSymbols)) // 63
)

func GenerateShortPath(originalUrl string, salt int) string {

	if salt > 0 {
		originalUrl = originalUrl + "#" + string(rune(salt))
	}

	hash := sha256.Sum256([]byte(originalUrl))

	// Переводим hash в число
	num := new(big.Int).SetBytes(hash[:])

	// Кодируем в путь для URL
	shortPath := make([]byte, maxLength)
	baseBig := big.NewInt(base)

	for i := maxLength - 1; i >= 0; i-- {
		mod := new(big.Int)
		num.DivMod(num, baseBig, mod)
		shortPath[i] = availableSymbols[mod.Int64()]
	}

	return string(shortPath)
}
