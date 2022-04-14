package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Metehan1994/final-project/internal/api"
	"github.com/Metehan1994/final-project/internal/category"
	httpErrors "github.com/Metehan1994/final-project/internal/httpErrors"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/pkg/config"
	jwtHelper "github.com/Metehan1994/final-project/pkg/jwt"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type userHandler struct {
	cfg          *config.Config
	userRepo     *UserRepository
	categoryRepo *category.CategoryRepository
	productRepo  *product.ProductRepository
}

func NewUserHandler(r *gin.RouterGroup, cfg *config.Config, userRepo *UserRepository, categoryRepo *category.CategoryRepository,
	productRepo *product.ProductRepository) {
	uHandler := userHandler{
		cfg:          cfg,
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
	r.POST("/login", uHandler.login)
	r.POST("/signup", uHandler.signUp)
	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	//r.POST("/AddToCart", uHandler.AddToCart)

	r.Use(mw.AuthMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/admin/addBulkCategory", uHandler.addBulkCategory)
	r.POST("/decode", uHandler.VerifyToken)
	r.POST("/admin/createProduct", uHandler.createProduct)
	r.PUT("/admin/updateProduct/:id", uHandler.updateProduct)
	r.DELETE("/admin/deleteProduct/:id", uHandler.deleteProduct)

}

func (u *userHandler) login(c *gin.Context) {
	var req api.Login
	if err := c.Bind(&req); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}
	user := u.userRepo.GetUserByEmail(*req.Email)
	fmt.Println(user)
	if user.Email == "" {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "user not found", nil)))
		return
	}
	err := ComparePasswordWithHashedOne(user.Password, *req.Password)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "password is wrong", nil)))
		return
	}
	apiUser := userToResponse(user)
	roles := RoleConvertToStringSlice(apiUser.IsAdmin)
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": apiUser.Username,
		"email":    apiUser.Email,
		"userId":   apiUser.ID,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"roles":    roles,
	})
	token := jwtHelper.GenerateToken(jwtClaims, u.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, token)
}

func (u *userHandler) signUp(c *gin.Context) {
	var req api.SignUp
	if err := c.Bind(&req); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "check your request body", nil)))
		return
	}
	if !ParseEmail(*req.Email) {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "The email is not in accepted format.", nil)))
		return
	}
	user := u.userRepo.GetUserByEmail(*req.Email)
	if user.Email != "" {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "The email is available on the system. Use it to connect.", nil)))
		return
	}
	user = u.userRepo.GetUserByUsername(*req.Username)
	if user.Username != "" {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "The username is taken by someone else. Try another one.", nil)))
		return
	}
	if *req.Password != *req.PasswordConfirm {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "Your password is not compatible with second entry.", nil)))
		return
	}
	password, _ := GenerateHashedPass(*req.Password)
	*req.Password = password
	DBUser := signedUpUserToDBUser(&req)
	user, err := u.userRepo.CreateNewUser(DBUser)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusInternalServerError, "Internal Error.", nil)))
		return
	}
	apiUser := userToResponse(user)
	roles := RoleConvertToStringSlice(apiUser.IsAdmin)
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": apiUser.Username,
		"email":    apiUser.Email,
		"userId":   apiUser.ID,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"roles":    roles,
	})
	token := jwtHelper.GenerateToken(jwtClaims, u.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, token)
}

func (u *userHandler) addBulkCategory(c *gin.Context) {
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
	ReadCSVforCategory(fileDir, u.categoryRepo)
}

func (u *userHandler) VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, decodedClaims)
}

func (u *userHandler) deleteProduct(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("id"))
	err := u.productRepo.DeleteById(idint)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, "The product is successfully deleted.")
}

func (u *userHandler) createProduct(c *gin.Context) {
	productBody := &api.Product{}
	if err := c.Bind(&productBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	productRepo, err := u.productRepo.Create(*product.ResponseToProduct(productBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, product.ProductToResponse(productRepo))
}

func (u *userHandler) updateProduct(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("id"))

	productBody2, err := u.productRepo.GetByID(idint)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Bind(&productBody2); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	prod, err := u.productRepo.Update(*productBody2)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, product.ProductToResponse(prod))
}
