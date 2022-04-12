package product

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
)

func ProductToResponse(p *models.Product) *api.Product {
	price := int64(p.Price)
	name := p.Name
	stockCode := p.StockCode
	return &api.Product{
		Name:        &name,
		Description: p.Description,
		Price:       &price,
		Quantity:    int64(p.Quantity),
		StockCode:   &stockCode,
	}
}
