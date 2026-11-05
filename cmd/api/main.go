package main

import (
	"log"
	"orders/internal/common/middleware"
	"orders/internal/domain/auth"
	"orders/internal/domain/inventory"
	inventoryRepository "orders/internal/domain/inventory/repository"
	"orders/internal/domain/product"
	productRepository "orders/internal/domain/product/repository"
	"orders/internal/domain/user"
	userRepository "orders/internal/domain/user/repository"
	"orders/internal/infrasctructure/database"
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

	authMiddleware := middleware.AuthMiddleware(authService)

	inventoryRepo := inventoryRepository.NewGormInventoryRepository(db)
	inventoryService := inventory.NewInventoryService(inventoryRepo)
	inventoryHandler := inventory.NewInventoryHandler(inventoryService)
	inventoryHandler.RegisterRoutes(api)

	productRepo := productRepository.NewGormProductRepository(db)
	productService := product.NewProductService(productRepo, inventoryService)
	productHandler := product.NewProductHandler(productService, authService)
	productHandler.RegisterRoutes(api, authMiddleware)

	if err := router.Run(":3000"); err != nil {
		panic(err)
	}
}
