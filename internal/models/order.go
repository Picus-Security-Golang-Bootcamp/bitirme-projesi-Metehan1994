package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderStatus int       `gorm:"default:0"`

	Items []OrderItem `gorm:"foreignKey:OrderID"`

	CartID uuid.UUID
	Cart   Cart `gorm:"foreignKey:CartID"`

	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	UserID     uuid.UUID `gorm:"OnDelete:SET NULL"`
	TotalPrice int
}

type OrderItem struct {
	gorm.Model
	Order   *Order    `json:"cart,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	OrderID uuid.UUID `gorm:"not null"`

	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"OnDelete:SET NULL"`
	Price     int     `gorm:"not null"`
}

func (order *Order) GetOrderStatusAsString() string {
	switch order.OrderStatus {
	case 0:
		return "open"
	case 1:
		return "completed"
	case 2:
		return "canceled"
	default:
		return "unknown"
	}
}
