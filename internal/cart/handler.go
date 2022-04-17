package cart

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	jwt_helper "github.com/Metehan1994/final-project/pkg/jwt"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cfg         *config.Config
	userRepo    *user.UserRepository
	cartService *CartService
}

func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, userRepo *user.UserRepository, cartService *CartService) {
	cart := &CartHandler{
		userRepo:    userRepo,
		cfg:         cfg,
		cartService: cartService,
	}

	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/addToCart/productId/:id/quantity/:quantity", cart.AddToCart)
	r.GET("/listCartItems", cart.ListCartItems)
	r.DELETE("/deleteItem/:itemId", cart.DeleteItem)
	r.PUT("/updateItem/:itemId/quantity/:quantity", cart.UpdateQuantity)
}

func (cHandler *CartHandler) AddToCart(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	//userDB := cHandler.userRepo.GetUserByEmail(user.Email)
	userDB, err := cHandler.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cart, Info := cHandler.cartService.GetOrCreateCart(userDB.ID)
	c.JSON(http.StatusOK, Info)
	fmt.Println(cart)
	idint, _ := strconv.Atoi(c.Param("id"))
	quantityint, _ := strconv.Atoi(c.Param("quantity"))
	product := cHandler.cartService.GetProductByID(idint)
	// if err != nil {
	// 	zap.L().Info(err.Error())
	// }
	Updatedcart, err := cHandler.cartService.AddItem(cart, product, quantityint)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	} else {
		cHandler.cartService.UpdateCartInDB(Updatedcart)
		cartBody := CartToResponse(Updatedcart)
		c.JSON(http.StatusAccepted, cartBody)
	}
}

func (cHandler *CartHandler) ListCartItems(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)

	userDB, err := cHandler.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cart := cHandler.cartService.GetCartByUserID(userDB.ID)
	cartBody := CartToResponse(cart)
	c.JSON(http.StatusOK, cartBody)
}

func (cHandler *CartHandler) DeleteItem(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)

	userDB, err := cHandler.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	idint, _ := strconv.Atoi(c.Param("itemId"))
	err = cHandler.cartService.DeleteItem(userDB.ID, idint)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "The product is successfully deleted.")
}

func (cHandler *CartHandler) UpdateQuantity(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("itemId"))
	quantity, _ := strconv.Atoi(c.Param("quantity"))
	user := c.MustGet("user").(*jwt_helper.DecodedToken)

	userDB, err := cHandler.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	err = cHandler.cartService.UpdateQuantityById(userDB.ID, idint, quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "The quantity is successfully updated.")
	cartUpdated := cHandler.cartService.GetCartByUserID(userDB.ID)
	cartBody := CartToResponse(cartUpdated)
	c.JSON(http.StatusOK, cartBody)
}
