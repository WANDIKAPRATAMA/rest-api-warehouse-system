package configs

import (
	controller "auth-service/internal/controllers"
	middleware "auth-service/internal/middlewares"
	"auth-service/internal/repositorys"
	route "auth-service/internal/routes"
	usecase "auth-service/internal/usecases"
	"auth-service/internal/utils"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validate    *validator.Validate
	Viper       *viper.Viper
	RedisClient *redis.Client
}

// AppConfig fungsi untuk setup app
func NewAppConfig(config *AppConfig) {

	jwtUtils := utils.NewJWTCfg(config.Viper, config.RedisClient)

	userRepo := repositorys.NewUserRepository(config.DB, config.Log)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.Log, config.Validate, config.Viper, jwtUtils)
	authController := controller.NewAuthController(authUseCase, config.Log, config.Validate)
	authMiddleware := middleware.NewAuth(authUseCase, config.Log, config.Viper, jwtUtils)

	productRepo := repositorys.NewProductRepository(config.DB)
	productUseCase := usecase.NewProductUseCase(productRepo, config.Log, config.Validate)
	productController := controller.NewProductController(productUseCase, config.Log, config.Validate)
	productMiddleware := middleware.NewProductMiddleware(productUseCase, config.Log)

	authRoutesConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}

	productRouteConfig := route.ProductRouteConfig{
		App:               config.App,
		ProductController: productController,
		ProductMiddleware: productMiddleware,
		AuthMiddleware:    authMiddleware,
	}

	productRouteConfig.Setup()
	authRoutesConfig.Setup()

	config.Log.Info("Server starting on :8080")
	if err := config.App.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
