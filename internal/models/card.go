package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	CartStatus int `gorm:"default:0"`

	CartItems []CartItem `gorm:"foreignKey:CartID"`

	User           User `gorm:"foreignKey:UserID:"`
	UserID         int
	CartItemsCount int `gorm:"-"`
	TotalPrice     int
}

type CartItem struct {
	gorm.Model
	Cart   Cart
	CartID uint `gorm:"not null"`

	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Amount    int     //`gorm:"not null"`
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
