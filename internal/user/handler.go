package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Metehan1994/final-project/internal/api"
	httpErrors "github.com/Metehan1994/final-project/internal/httpErrors"
	"github.com/Metehan1994/final-project/pkg/config"
	jwtHelper "github.com/Metehan1994/final-project/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type userHandler struct {
	cfg      *config.Config
	userRepo *UserRepository
}

func NewUserHandler(r *gin.RouterGroup, cfg *config.Config, userRepo *UserRepository) {
	uHandler := userHandler{
		cfg:      cfg,
		userRepo: userRepo,
	}
	r.POST("/login", uHandler.login)
	r.POST("/signup", uHandler.signUp)

	//r.Use(mw.AuthMiddleware(cfg.ServerConfig.RoutePrefix))
	//r.POST("/decode", a.VerifyToken)

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
	err := ComparePassWithHashed(user.Password, *req.Password)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.NewRestError(http.StatusBadRequest, "password is wrong", nil)))
		return
	}
	apiUser := userToResponse(user)
	roles := RoleConvertToStringSlice(apiUser.IsAdmin)
	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": apiUser.ID,
		"email":  apiUser.Email,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  roles,
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
		"userId": apiUser.ID,
		"email":  apiUser.Email,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
		"roles":  roles,
	})
	token := jwtHelper.GenerateToken(jwtClaims, u.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, token)
}

func (u *userHandler) VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	decodedClaims := jwtHelper.VerifyToken(token, u.cfg.JWTConfig.SecretKey)
	c.JSON(http.StatusOK, decodedClaims)
}
