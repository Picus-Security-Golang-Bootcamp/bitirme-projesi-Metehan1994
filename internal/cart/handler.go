package cart

import (
	"net/http"
	"strconv"

	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	jwt_helper "github.com/Metehan1994/final-project/pkg/jwt"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cfg          *config.Config
	Cartrepo     *CartRepository
	productRepo  *product.ProductRepository
	cartItemRepo *CartItemRepository
	userRepo     *user.UserRepository
}

func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, Cartrepo *CartRepository, productRepo *product.ProductRepository,
	cartItemRepo *CartItemRepository, userRepo *user.UserRepository) {
	cart := &CartHandler{
		Cartrepo:     Cartrepo,
		productRepo:  productRepo,
		cartItemRepo: cartItemRepo,
		userRepo:     userRepo,
		cfg:          cfg,
	}

	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/addToCart/productId/:id/quantity/:quantity", cart.AddToCart)
	r.GET("/listCartItems", cart.ListCartItems)
	r.DELETE("/deleteItem/:itemId", cart.DeleteItem)
	r.PUT("/updateItem/:itemId/quantity/:quantity", cart.UpdateQuantity)
}

func (cHandler *CartHandler) AddToCart(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := cHandler.userRepo.GetUserByEmail(user.Email)

	cart, Info := cHandler.Cartrepo.GetOrCreateCart(userDB.ID)
	c.JSON(http.StatusOK, Info)

	idint, _ := strconv.Atoi(c.Param("id"))
	quantityint, _ := strconv.Atoi(c.Param("quantity"))
	product, _ := cHandler.productRepo.GetByID(idint)
	// if err != nil {
	// 	zap.L().Info(err.Error())
	// }
	Updatedcart, err := cHandler.cartItemRepo.AddItem(cart, product, quantityint)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		cHandler.Cartrepo.Update(Updatedcart)
		cartBody := CartToResponse(Updatedcart)
		c.JSON(http.StatusAccepted, cartBody)
	}
}

func (cHandler *CartHandler) ListCartItems(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := cHandler.userRepo.GetUserByEmail(user.Email)

	cart := cHandler.Cartrepo.GetCartByUserID(userDB.ID)
	cartBody := CartToResponse(cart)
	c.JSON(http.StatusOK, cartBody)
}

func (cHandler *CartHandler) DeleteItem(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := cHandler.userRepo.GetUserByEmail(user.Email)
	cart := cHandler.Cartrepo.GetCartByUserID(userDB.ID)
	idint, _ := strconv.Atoi(c.Param("itemId"))
	err := cHandler.Cartrepo.DeleteItemByID(cart, idint)
	//err := cHandler.cartItemRepo.DeleteById(uint(idint))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, "The product is successfully deleted.")
}

func (cHandler *CartHandler) UpdateQuantity(c *gin.Context) {
	idint, _ := strconv.Atoi(c.Param("itemId"))
	quantity, _ := strconv.Atoi(c.Param("quantity"))
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := cHandler.userRepo.GetUserByEmail(user.Email)
	cart := cHandler.Cartrepo.GetCartByUserID(userDB.ID)
	err := cHandler.Cartrepo.UpdateQuantityById(cart, idint, quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "The quantity is successfully updated.")
	cartUpdated := cHandler.Cartrepo.GetCartByUserID(userDB.ID)
	cartBody := CartToResponse(cartUpdated)
	c.JSON(http.StatusOK, cartBody)
}
