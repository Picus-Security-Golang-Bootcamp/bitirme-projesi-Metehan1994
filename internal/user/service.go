package user

import (
	"net/mail"
	"time"

	"github.com/Metehan1994/final-project/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

//GenerateHashedPass generates a salted hash for the input string
func GenerateHashedPass(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

//ComparePasswordWithHashedOne compares string to generated hash
func ComparePasswordWithHashedOne(hash string, s string) error {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming)
}

//ParseEmail checks the email is in acceptable format or not
func ParseEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

//JWTClaimsGenerator generates a JWTClaims for a given user
func JWTClaimsGenerator(user *models.User) *jwt.Token {
	apiUser := userToResponse(user)
	roles := RoleConvertToStringSlice(apiUser.IsAdmin)
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": apiUser.Username,
		"email":    apiUser.Email,
		"userId":   apiUser.ID,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
		"roles":    roles,
	})
	return jwtClaims
}
