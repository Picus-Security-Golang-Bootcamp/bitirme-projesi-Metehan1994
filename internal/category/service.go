package category

import (
	"encoding/csv"
	"os"

	"github.com/Metehan1994/final-project/internal/models"
	"go.uber.org/zap"
)

//ReadCSV reads csv and returns products and categories
func ReadCSVforCategory(filename string, categoryRepo *CategoryRepository) {
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
