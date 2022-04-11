package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string
	LastName  string
	Username  string
	Email     string `json:"email" gorm:"unique"`
	Password  string
	IsAdmin   bool
}

func (User) TableName() string {
	//default table name
	return "users"
}
