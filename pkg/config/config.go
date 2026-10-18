package config

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	DatePattern = "20060102"
)

func GetPortFromEnv() int {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return 7540
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return 7540
	}
	return intPort
}

func GetDBFromEnv() string {
	db := os.Getenv("TODO_DBFILE")
	if db == "" {
		return "scheduler.db"
	}
	return db
}

func GetPasswordFromEnv() string {
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		return "qwerty12345"
	}
	return password
}

func GetToken(password string) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(password))
}
