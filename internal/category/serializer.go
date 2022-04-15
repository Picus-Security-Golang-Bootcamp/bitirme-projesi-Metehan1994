package category

import (
	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/models"
	"github.com/Metehan1994/final-project/internal/product"
)

func CategoryToResponse(c *models.Category) *api.Category {
	products := make([]*api.Product, 0)
	for _, prod := range c.Products {
		products = append(products, product.ProductToResponse(&prod))
	}

	name := c.Name
	return &api.Category{
		Name:        &name,
		Description: c.Description,
		Products:    products,
	}
}

func CategoriesToResponse(c []*models.Category) []*api.Category {
	categories := make([]*api.Category, 0)
	for _, b := range c {
		categories = append(categories, CategoryToResponse(b))
	}
	return categories
}
