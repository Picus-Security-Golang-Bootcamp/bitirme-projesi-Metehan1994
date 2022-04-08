package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"unique"`
	LastName  string
	Username  string
	Email     string
	Password  string
	IsAdmin   bool
	IsUser    bool
}

func (User) TableName() string {
	//default table name
	return "users"
}
