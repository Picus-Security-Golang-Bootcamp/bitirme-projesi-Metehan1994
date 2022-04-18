package product

import (
	"net/http"
	"strconv"

	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/httpErrors"
	"github.com/Metehan1994/final-project/pkg/config"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/Metehan1994/final-project/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type ProductHandler struct {
	productRepo *ProductRepository
	cfg         *config.Config
}

func NewProductHandler(r *gin.RouterGroup, productRepo *ProductRepository, cfg *config.Config) {
	h := &ProductHandler{productRepo: productRepo, cfg: cfg}

	r.GET("/list", h.ProductList)
	r.GET("/search/name/:word", h.SearchByName)
	r.GET("/search/sku/:word", h.SearchBySku)

	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/admin/createProduct", h.createProduct)
	r.PUT("/admin/updateProduct/:id", h.updateProduct)
	r.DELETE("/admin/deleteProduct/:id", h.deleteProduct)
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

func (pro *ProductHandler) createProduct(c *gin.Context) {
	productBody := &api.Product{}
	if err := c.Bind(&productBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	productRepo, err := pro.productRepo.Create(*ResponseToProduct(productBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductToResponse(productRepo))
}

func (pro *ProductHandler) deleteProduct(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("id"))
	err := pro.productRepo.DeleteById(idint)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, "The product is successfully deleted.")
}

func (pro *ProductHandler) updateProduct(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("id"))

	productBody, err := pro.productRepo.GetByID(idint)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Bind(&productBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	prod, err := pro.productRepo.Update(*productBody)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductToResponse(prod))
}
