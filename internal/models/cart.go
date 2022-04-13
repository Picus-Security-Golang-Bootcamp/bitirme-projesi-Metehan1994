package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status int       `gorm:"default:0"`

	Items []CartItem

	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	UserID     uuid.UUID `gorm:"OnDelete:SET NULL"`
	TotalPrice int
}

type CartItem struct {
	gorm.Model
	Cart   *Cart     `json:"cart,omitempty" gorm:"foreignKey:CartID;references:ID"`
	CartID uuid.UUID `gorm:"not null"`

	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"OnDelete:SET NULL"`
	Price     int     `gorm:"not null"`
	Amount    int     `gorm:"not null"`
}

func (cart *Cart) GetCartStatusAsString() string {
	switch cart.Status {
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
