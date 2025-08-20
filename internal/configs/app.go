package configs

import (
	controller "auth-service/internal/controllers"
	middleware "auth-service/internal/middlewares"
	"auth-service/internal/repositorys"
	route "auth-service/internal/routes"
	usecase "auth-service/internal/usecases"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AppConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Viper    *viper.Viper
}

// AppConfig fungsi untuk setup app
func NewAppConfig(config *AppConfig) {

	userRepo := repositorys.NewUserRepository(config.DB, config.Log)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.Log, config.Validate, config.Viper)
	authController := controller.NewAuthController(authUseCase, config.Log, config.Validate)
	authMiddleware := middleware.NewAuth(authUseCase, config.Log, config.Viper)
	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
	config.Log.Info("Server starting on :3000")
	if err := config.App.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
