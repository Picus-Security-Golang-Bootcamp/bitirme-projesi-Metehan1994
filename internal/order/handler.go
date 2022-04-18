package order

import (
	"net/http"

	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	jwt_helper "github.com/Metehan1994/final-project/pkg/jwt"
	mw "github.com/Metehan1994/final-project/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	cfg          *config.Config
	userRepo     *user.UserRepository
	orderService *OrderService
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, userRepo *user.UserRepository, orderService *OrderService) {
	order := &OrderHandler{
		userRepo:     userRepo,
		cfg:          cfg,
		orderService: orderService,
	}
	r.Use(mw.TokenExpControlMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("/completeOrder", order.CompleteOrder)
	r.GET("/listOrder", order.ListOrder)
	r.PUT("/cancelOrder/orderId/:id", order.CancelOrder)
}

func (o *OrderHandler) CompleteOrder(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB, err := o.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cart := o.orderService.GetCartByUserID(userDB.ID)
	if cart.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no available cart"})
		return
	}
	order, err := o.orderService.CompleteOrder(cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderWithItems, err := o.orderService.OrderGetWithItems(order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.orderService.DeleteCart(cart)
	orderBody := OrderToResponse(orderWithItems)
	c.JSON(http.StatusAccepted, orderBody)
}

func (o *OrderHandler) ListOrder(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB, err := o.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	orderList, err := o.orderService.GetOrdersByUserID(userDB.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ordersBody := OrderListToResponse(orderList)
	c.JSON(http.StatusOK, ordersBody)
}

func (o *OrderHandler) CancelOrder(c *gin.Context) {
	user := c.MustGet("user").(*jwt_helper.DecodedToken)
	userDB, err := o.userRepo.GetUserByID(user.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	idString := (c.Param("id"))
	order, err := o.orderService.CancelOrder(userDB.ID, idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "order canceled")
	orderBody := OrderToResponse(order)
	c.JSON(http.StatusAccepted, orderBody)
}
