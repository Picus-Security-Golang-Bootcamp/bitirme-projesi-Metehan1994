package main

import (
	"fmt"
	"log"

	"github.com/Metehan1994/final-project/internal/category"
	"github.com/Metehan1994/final-project/internal/product"
	"github.com/Metehan1994/final-project/internal/user"
	"github.com/Metehan1994/final-project/pkg/config"
	csvReader "github.com/Metehan1994/final-project/pkg/csv"
	db "github.com/Metehan1994/final-project/pkg/database"
	logger "github.com/Metehan1994/final-project/pkg/logging"
)

func main() {
	log.Println("Book store service starting...")

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

	// Product Repository
	productRepo := product.NewProductRepository(DB)
	productRepo.Migration()

	// Category Repository
	categoryRepo := category.NewCategoryRepository(DB)
	categoryRepo.Migration()

	//User Repository
	userRepo := user.NewUserRepository(DB)
	userRepo.Migration()

	//Initialize products&categories
	csvReader.ReadCSVforProducts("./pkg/csv/files/products.csv", categoryRepo, productRepo)
	csvReader.ReadCSVforUsers("./pkg/csv/files/users.csv", userRepo)
}
