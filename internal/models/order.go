package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderStatus    int `gorm:"default:0"`
	TrackingNumber string

	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`

	User            User `gorm:"foreignKey:UserID:"`
	UserID          int
	OrderItemsCount int `gorm:"-"`
	TotalPrice      int `gorm:"-"`
}

type OrderItem struct {
	gorm.Model
	Cart    Cart
	OrderID uint `gorm:"not null"`

	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Amount    int     //`gorm:"not null"`
}

func (order *Order) GetOrderStatusAsString() string {
	switch order.OrderStatus {
	case 0:
		return "proccessed"
	case 1:
		return "delivered"
	case 2:
		return "shipped"
	default:
		return "unknown"
	}
}
