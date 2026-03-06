package utils

import (
	"crypto/sha256"
	"math/big"
	"url-shortener-ozon/internal/domain/entities"
)

const (
	maxLength = 10
	alphabet  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base      = int64(len(alphabet)) // 63
)

func GenerateShortURL(url entities.InOutURL, salt int) entities.InOutURL {

	input := url.URL
	if salt > 0 {
		input = input + "#" + string(rune(salt))
	}

	hash := sha256.Sum256([]byte(input))

	// Переводим hash в число
	num := new(big.Int).SetBytes(hash[:])

	// Кодируем в base63
	result := make([]byte, maxLength)
	baseBig := big.NewInt(base)

	for i := maxLength - 1; i >= 0; i-- {
		mod := new(big.Int)
		num.DivMod(num, baseBig, mod)
		result[i] = alphabet[mod.Int64()]
	}

	return entities.InOutURL{
		URL: string(result),
	}
}
