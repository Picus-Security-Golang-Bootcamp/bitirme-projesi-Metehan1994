package category

import (
	"fmt"
	"net/http"

	"github.com/Metehan1994/final-project/internal/httpErrors"
	"github.com/Metehan1994/final-project/pkg/config"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/Metehan1994/final-project/pkg/pagination"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	Catrepo *CategoryRepository
	cfg     *config.Config
}

func NewCategoryHandler(r *gin.RouterGroup, Catrepo *CategoryRepository, cfg *config.Config) {
	h := &CategoryHandler{Catrepo: Catrepo, cfg: cfg}

	r.GET("/list", h.CategoryList)

	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/admin/addBulkCategory", h.addBulkCategory)
}

func (cat *CategoryHandler) CategoryList(c *gin.Context) {
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(c)
	categories, count := cat.Catrepo.ListCategoriesWithProducts(pageIndex, pageSize)
	paginatedResult := pagination.NewFromGinRequest(c, count)
	paginatedResult.Items = CategoriesToResponse(categories)
	c.JSON(http.StatusOK, paginatedResult)
}

func (cat *CategoryHandler) addBulkCategory(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "Cannot upload file.", nil)))
		return
	}
	fileDir := "pkg/csv/files/saved/" + file.Filename
	err = c.SaveUploadedFile(file, fileDir)
	if err != nil {
		zap.L().Fatal(err.Error())
	}
	c.JSON(http.StatusOK, fmt.Sprintf("'%s' is uploaded!", file.Filename))
	ReadCSVforCategory(fileDir, cat.Catrepo)
}
