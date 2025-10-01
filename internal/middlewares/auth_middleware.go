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
	usecase          usecases.AuthUseCase
	log              *logrus.Logger
	config           *viper.Viper
	jwtUtils         *utils.JWTConfig
	rateLimiterUtils *utils.RateLimiterUtil
}

func NewAuth(usecase usecases.AuthUseCase, log *logrus.Logger, config *viper.Viper, jwtUtils *utils.JWTConfig, rateLimiterUtil *utils.RateLimiterUtil) *AuthMiddleware {
	return &AuthMiddleware{usecase: usecase, log: log, config: config,
		jwtUtils:         jwtUtils,
		rateLimiterUtils: rateLimiterUtil}
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

	m.log.Printf("token: %v", tokenString)
	token, err := m.jwtUtils.ValidateToken(c.Context(), tokenString, utils.AccessToken)
	if err != nil || !token.Valid {
		m.log.Printf("error: %v", err)
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	if allow := m.rateLimiterUtils.IsAllowed(c.Context(), tokenString); !allow {
		return fiber.NewError(fiber.StatusTooManyRequests, "Too many requests")
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

type LocalKeys struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

func GetLocalKeys(c *fiber.Ctx) *LocalKeys {
	userID, _ := c.Locals("userID").(uuid.UUID)
	email, _ := c.Locals("email").(string)
	role, _ := c.Locals("role").(string)
	return &LocalKeys{
		UserID: userID,
		Email:  email,
		Role:   role,
	}
}
