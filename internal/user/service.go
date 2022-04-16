package user

import (
	"encoding/csv"
	"net/mail"
	"os"
	"time"

	"github.com/Metehan1994/final-project/internal/category"
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
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

//ReadCSV reads csv and returns products and categories
func ReadCSVforCategory(filename string, categoryRepo *category.CategoryRepository) {
	f, err := os.Open(filename)
	if err != nil {
		zap.L().Fatal("Cannot open csv file")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		zap.L().Fatal("Cannot read csv")
	}

	for _, line := range records[1:] {
		category := models.Category{}
		category.Name = line[0]
		category.Description = line[1]
		categoryRepo.InsertSampleData(&category)
	}
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
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"roles":    roles,
	})
	return jwtClaims
}
