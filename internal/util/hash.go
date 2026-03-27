package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, length int) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), length)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
