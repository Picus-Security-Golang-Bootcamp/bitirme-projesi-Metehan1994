package category

import (
	"net/http"

	"github.com/Metehan1994/final-project/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	Catrepo *CategoryRepository
}

func NewCategoryHandler(r *gin.RouterGroup, Catrepo *CategoryRepository) {
	h := &CategoryHandler{Catrepo: Catrepo}

	r.GET("/list", h.CategoryList)
}

func (cat *CategoryHandler) CategoryList(c *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)
	categories, count := cat.Catrepo.ListCategoriesWithProducts(pageIndex, pageSize)
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = CategoriesToResponse(categories)
	c.JSON(http.StatusOK, paginatedResult)
}
