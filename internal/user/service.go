package user

import (
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

//Generate a salted hash for the input string
func GenerateHashedPass(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//Compare string to generated hash
func ComparePassWithHashed(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

func ParseEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
