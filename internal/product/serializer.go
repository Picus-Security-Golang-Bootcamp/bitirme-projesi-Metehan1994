package product

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
)

func ProductToResponse(p *models.Product) *api.Product {
	name := p.Name
	stockCode := p.Sku
	CategoryName := p.Category.Name
	return &api.Product{
		Name:        &name,
		Description: p.Description,
		Price:       int64(p.Price),
		Quantity:    int64(p.Quantity),
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
	return &models.Product{
		Name:        name,
		Description: apiPro.Description,
		Price:       int(apiPro.Price),
		Quantity:    int(apiPro.Quantity),
		Sku:         stockCode,
		CategoryID:  uint(apiPro.Category.ID),
	}
}
