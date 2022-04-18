package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Metehan1994/final-project/internal/cart"
	"github.com/Metehan1994/final-project/internal/category"
	"github.com/Metehan1994/final-project/internal/order"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	csvReader "github.com/Metehan1994/final-project/pkg/csv"
	db "github.com/Metehan1994/final-project/pkg/database"
	"github.com/Metehan1994/final-project/pkg/graceful"
	logger "github.com/Metehan1994/final-project/pkg/logging"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Basket service starting...")

	// Set envs for local development
	cfg, err := config.LoadConfig("./pkg/config/config-local")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	// Set global logger
	logger.NewLogger(cfg)
	defer logger.Close()

	// Connect DB
	DB := db.Connect(cfg)
	fmt.Println("Connected to DB")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs * int64(time.Second)),
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs * int64(time.Second)),
	}

	// Router group
	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix)
	authRouter := rootRouter.Group("/user")
	categoryRouter := rootRouter.Group("/category")
	productRouter := rootRouter.Group("/product")
	cartRouter := rootRouter.Group("/cart")
	orderRouter := rootRouter.Group("/order")

	// Product Repository
	productRepo := product.NewProductRepository(DB)
	productRepo.Migration()
	product.NewProductHandler(productRouter, productRepo, cfg)

	// Category Repository
	categoryRepo := category.NewCategoryRepository(DB)
	categoryRepo.Migration()
	category.NewCategoryHandler(categoryRouter, categoryRepo, cfg)

	//User Repository
	userRepo := user.NewUserRepository(DB)
	userRepo.Migration()
	user.NewUserHandler(authRouter, cfg, userRepo)

	//Cart Repository
	cartItemRepo := cart.NewCartItemRepository(DB)
	cartRepo := cart.NewCartRepository(DB, cartItemRepo)
	cartRepo.Migration()
	cartItemRepo.Migration()
	cartService := cart.InitializeCartService(cartRepo, productRepo, cartItemRepo)
	cart.NewCartHandler(cartRouter, cfg, userRepo, cartService)

	//Order Repository
	orderItemRepo := order.NewOrderItemRepository(DB)
	orderRepo := order.NewOrderRepository(DB)
	orderRepo.Migration()
	orderItemRepo.Migration()
	orderService := order.InitializeOrderService(orderRepo, productRepo, orderItemRepo, cartRepo)
	order.NewOrderHandler(orderRouter, cfg, userRepo, orderService)

	//Initialize products&categories&users
	csvReader.ReadCSVforProducts("./pkg/csv/files/products.csv", categoryRepo, productRepo)
	csvReader.ReadCSVforUsers("./pkg/csv/files/users.csv", userRepo)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Basket service started")
	graceful.ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))
}
