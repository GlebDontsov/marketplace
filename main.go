package main

import (
	_ "marketplace/docs"
	"marketplace/internal/config"
	"marketplace/internal/controller"
	"marketplace/internal/middleware"
	"marketplace/internal/repository"
	"marketplace/internal/service"
	"marketplace/pkg/database"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := database.InitDB(database.DBConfig{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		DBName:   cfg.DBConfig.DBName,
		SSLMode:  cfg.DBConfig.SSLMode,
	})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&repository.User{}, &repository.Advertisement{}); err != nil {
		panic("failed to migrate database")
	}

	userRepo := repository.NewUserRepository(db)
	adRepo := repository.NewAdvertisementRepository(db)

	userService := service.NewUserService(userRepo)
	adService := service.NewAdvertisementService(adRepo)

	userController := controller.NewUserController(*userService)
	adController := controller.NewAdvertisementController(*adService)

	r := gin.Default()

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	adRoutes := r.Group("/ads")
	adRoutes.Use(middleware.AuthMiddleware())
	{
		adRoutes.POST("/", adController.CreateAd)
		adRoutes.GET("/", adController.GetAds)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
