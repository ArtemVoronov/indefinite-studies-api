package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func CreateSHA512HashHexEncoded(str string) string {
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateSHA256HashHexEncoded(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Contains(array []string, elem string) bool {
	for _, n := range array {
		if elem == n {
			return true
		}
	}
	return false
}
