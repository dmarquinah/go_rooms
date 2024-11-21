package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GenerateRandomCode(max_digits int) (string, error) {
	buffer := make([]byte, max_digits)
	for i := 0; i < max_digits; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(dictionary))))
		if err != nil {
			return "", errors.New("unable to generate code")
		}
		buffer[i] = dictionary[num.Int64()]
	}

	return string(buffer), nil
}

const dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
