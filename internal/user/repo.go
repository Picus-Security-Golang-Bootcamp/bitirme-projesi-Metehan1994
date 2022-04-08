package user

import (
	"github.com/Metehan1994/final-project/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (c *UserRepository) Migration() {
	c.db.AutoMigrate(&models.User{})
}

func (c *UserRepository) InsertSampleData(user *models.User) {
	result := c.db.Unscoped().Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	if result.Error != nil {
		zap.L().Fatal("Cannot insert data into DB") //Check Error
	}
}
