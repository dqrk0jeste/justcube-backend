package util

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func VerifyPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
