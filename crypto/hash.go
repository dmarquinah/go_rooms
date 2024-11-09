package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Setting up a reference cost value
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(bytes), err
}

func VerifyPassword(bodyPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(bodyPassword))
	return err == nil
}
