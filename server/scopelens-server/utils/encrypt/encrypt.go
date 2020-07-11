package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// Encrypt password and return hashedPassword and error.
func PasswordEncrypt(password string) (string, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func PasswordVerification(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
