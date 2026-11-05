package main

import (
	"log"
	"orders/internal/auth"
	"orders/internal/infrasctructure/database"
	"orders/internal/inventory"
	inventoryRepository "orders/internal/inventory/repository"
	"orders/internal/product"
	productRepository "orders/internal/product/repository"
	"orders/internal/user"
	userRepository "orders/internal/user/repository"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca dari sistem env")
	}

	db, err := database.NewPostgres(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	router := gin.Default()

	api := router.Group("/api")

	userRepo := userRepository.NewGormUserRepository(db)
	userService := user.NewUserService(userRepo)

	authService := auth.NewAuthService(
		os.Getenv("JWT_SECRET"),
		24,
		userService,
	)
	authHandler := auth.NewAuthHandler(authService)
	authHandler.RegisterRoutes(api)

	inventoryRepo := inventoryRepository.NewGormInventoryRepository(db)
	inventoryService := inventory.NewInventoryService(inventoryRepo)
	inventoryHandler := inventory.NewInventoryHandler(inventoryService)
	inventoryHandler.RegisterRoutes(api)

	productRepo := productRepository.NewGormProductRepository(db)
	productService := product.NewProductService(productRepo, inventoryService)
	productHandler := product.NewProductHandler(productService, authService)
	productHandler.RegisterRoutes(api)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
