package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	CartStatus int       `gorm:"default:0"`

	CartItems []CartItem

	User           User `gorm:"foreignKey:UserID"`
	UserID         uuid.UUID
	CartItemsCount int `gorm:"-"`
	TotalPrice     int
}

type CartItem struct {
	gorm.Model
	Cart   Cart      `gorm:"foreignKey:CartID;references:ID"`
	CartID uuid.UUID `gorm:"not null"`

	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Price     int     `gorm:"not null"`
	Amount    int     `gorm:"not null"`
}

func (cart *Cart) GetCartStatusAsString() string {
	switch cart.CartStatus {
	case 0:
		return "open"
	case 1:
		return "submitted"
	case 2:
		return "deleted"
	default:
		return "unknown"
	}
}
