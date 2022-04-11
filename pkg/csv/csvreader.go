package csvReader

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/Metehan1994/final-project/internal/category"
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/internal/user"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//ReadCSV reads csv and returns products and categories
func ReadCSVforProducts(filename string, categoryRepo *category.CategoryRepository, productRepo *product.ProductRepository) {
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
		product := models.Product{}
		product.Name = line[0]
		product.Price, _ = strconv.Atoi(line[1])
		product.Quantity, _ = strconv.Atoi(line[2])
		product.StockCode = line[3]
		product.Category.Name = line[4]
		newCategory := categoryRepo.InsertSampleData(&product.Category)
		product.CategoryID = newCategory.ID
		productRepo.InsertSampleData(&product)
	}
}

func ReadCSVforUsers(filename string, userRepo *user.UserRepository) {
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
		user := models.User{}
		user.FirstName = line[0]
		user.LastName = line[1]
		user.Username = line[2]
		user.Email = line[3]
		user.Password = line[4]
		user.IsAdmin, _ = strconv.ParseBool(line[5])

		DBUser := userRepo.GetUserByEmail(user.Email)
		if user.Email == DBUser.Email {
			user.ID = DBUser.ID
		} else {
			user.ID = uuid.New()
		}
		//user.ID = uuid.New()
		userRepo.InsertSampleData(&user)
	}
}
