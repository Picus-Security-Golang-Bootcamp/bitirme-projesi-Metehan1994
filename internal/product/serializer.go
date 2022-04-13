package product

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
)

func ProductToResponse(p *models.Product) *api.Product {
	name := p.Name
	stockCode := p.Sku
	CategoryName := p.Category.Name
	price := int64(p.Price)
	quantity := int64(p.Quantity)
	return &api.Product{
		Name:        &name,
		Description: p.Description,
		Price:       &price,
		Quantity:    &quantity,
		Sku:         &stockCode,
		Category: &api.CategoryWithoutRequiredName{
			ID:   int64(p.Category.ID),
			Name: CategoryName,
		},
	}
}

func ResponseToProduct(apiPro *api.Product) *models.Product {
	name := *apiPro.Name
	stockCode := *apiPro.Sku
	price := *apiPro.Price
	quantity := *apiPro.Quantity
	return &models.Product{
		Name:        name,
		Description: apiPro.Description,
		Price:       int(price),
		Quantity:    int(quantity),
		Sku:         stockCode,
		CategoryID:  uint(apiPro.Category.ID),
	}
}
