package middleware

import (
	"auth-service/internal/usecases"
	"auth-service/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	usecase usecases.AuthUseCase
	log     *logrus.Logger
	config  *viper.Viper
}

func NewAuth(usecase usecases.AuthUseCase, log *logrus.Logger, config *viper.Viper) *AuthMiddleware {
	return &AuthMiddleware{usecase: usecase, log: log, config: config}
}

func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}
	secret := m.config.GetString("jwt.secret")
	token, err := utils.ValidateToken(secret, tokenString)
	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
	}

	c.Locals("userID", userID)
	c.Locals("email", claims["email"].(string))
	c.Locals("role", claims["role"].(string))

	return c.Next()
}
