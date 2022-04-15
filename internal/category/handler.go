package category

import (
	"net/http"

	httpErrors "github.com/Metehan1994/final-project/internal/httpErrors"
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
	categories, err := cat.Catrepo.ListCategoriesWithProducts()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, CategoriesToResponse(categories))
}
