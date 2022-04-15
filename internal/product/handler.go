package product

import (
	"net/http"

	"github.com/Metehan1994/final-project/internal/httpErrors"
	"github.com/Metehan1994/final-project/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo *ProductRepository
}

func NewProductHandler(r *gin.RouterGroup, productRepo *ProductRepository) {
	h := &ProductHandler{productRepo: productRepo}

	r.GET("/list", h.ProductList)
	r.GET("/search/name/:word", h.SearchByName)
	r.GET("/search/sku/:word", h.SearchBySku)
}

func (pro *ProductHandler) ProductList(c *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)
	products, count := pro.productRepo.ListProductsWithCategories(pageIndex, pageSize)
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = ProductListToResponse(products)
	c.JSON(http.StatusOK, paginatedResult)
}

func (pro *ProductHandler) SearchByName(c *gin.Context) {
	products, err := pro.productRepo.SearchByName(c.Param("word"))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductListToResponse(products))
}

func (pro *ProductHandler) SearchBySku(c *gin.Context) {
	products, err := pro.productRepo.SearchBySku(c.Param("word"))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductListToResponse(products))
}
