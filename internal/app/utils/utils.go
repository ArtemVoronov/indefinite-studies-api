package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"log"
	"os"
	"strconv"
	"time"
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

func EnvVar(varName string) string {
	val, valExists := os.LookupEnv(varName)
	if !valExists {
		log.Fatalf("Missed enviroment variable: %s. Check the .env file or OS enviroment vars", varName)
	}
	return val
}

func EnvVarDefault(varName string, defaultValue string) string {
	val, valExists := os.LookupEnv(varName)
	if !valExists {
		return defaultValue
	}
	return val
}

func EnvVarInt(varName string) int {
	val := EnvVar(varName)
	result, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Wrong value of environment variable: %s. It should be integer number", varName)
	}
	return result
}

func EnvVarDuration(varName string, unit time.Duration) time.Duration {
	val := EnvVarInt(varName)
	return unit * time.Duration(val)
}

func EnvVarBytes(varName string) []byte {
	val := EnvVar(varName)
	return []byte(val)
}
