package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
	crypted, err := bcrypt.GenerateFromPassword([]byte(password), 3)
	return string(crypted), err
}

func CorrectPassword(crypted string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(crypted), []byte(password))
	return err == nil
}
