package convert_functions

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	ShortURLLength = 10
	AllowedChars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func Convert_url_to_string(url string) (string, error) {
	allowedCharsLength := big.NewInt(int64(len(AllowedChars)))
	str := make([]byte, ShortURLLength)

	for i := 0; i < ShortURLLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, allowedCharsLength)
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		str[i] = AllowedChars[randomIndex.Int64()]

	}
	return string(str), nil
}

func Contains(arr []byte, char byte) bool {
	for _, c := range arr {
		if c == char {
			return true
		}
	}
	return false
}
