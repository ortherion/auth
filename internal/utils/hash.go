package utils

import (
	"crypto"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func GeneratePasswordHash(password string) string {
	hash := crypto.SHA1.New()
	hash.Write([]byte(password))
	return string(hash.Sum([]byte(os.Getenv("PASSWORD_SALT"))))
	//return string(hash.Sum([]byte("")))
}

func HashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswordAndHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
