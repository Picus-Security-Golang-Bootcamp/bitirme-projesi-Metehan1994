package order

import (
	"net/http"

	"github.com/Metehan1994/final-project/internal/cart"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	jwt_helper "github.com/Metehan1994/final-project/pkg/jwt"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	cfg           *config.Config
	orderRepo     *OrderRepository
	productRepo   *product.ProductRepository
	orderItemRepo *OrderItemRepository
	userRepo      *user.UserRepository
	cartRepo      *cart.CartRepository
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, orderRepo *OrderRepository, productRepo *product.ProductRepository,
	orderItemRepo *OrderItemRepository, userRepo *user.UserRepository, cartRepo *cart.CartRepository) {
	order := &OrderHandler{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		orderItemRepo: orderItemRepo,
		userRepo:      userRepo,
		cfg:           cfg,
		cartRepo:      cartRepo,
	}
	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/completeOrder", order.CompleteOrder)
	r.GET("/listOrder", order.ListOrder)
	r.PUT("/cancelOrder", order.ListOrder)
}

func (o *OrderHandler) CompleteOrder(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := o.userRepo.GetUserByEmail(user.Email)
	cart := o.cartRepo.GetCartByUserID(userDB.ID)
	if cart.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, "no avaiable cart")
		return
	}
	order, err := o.orderRepo.CompleteOrder(cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	orderWithItems, err := o.orderRepo.OrderGetWithItems(order)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	o.cartRepo.DeleteCart(cart)
	orderBody := OrderToResponse(orderWithItems)
	c.JSON(http.StatusAccepted, orderBody)
}

func (o *OrderHandler) ListOrder(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB := o.userRepo.GetUserByEmail(user.Email)

	orderList, err := o.orderRepo.GetOrdersByUserID(userDB.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ordersBody := OrderListToResponse(orderList)
	c.JSON(http.StatusOK, ordersBody)
}
