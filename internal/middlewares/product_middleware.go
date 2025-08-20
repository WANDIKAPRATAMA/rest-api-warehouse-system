package middleware

import (
	"auth-service/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProductMiddleware struct {
	usecase usecases.ProductUseCase
	log     *logrus.Logger
}

func NewProductMiddleware(usecase usecases.ProductUseCase, log *logrus.Logger) *ProductMiddleware {
	return &ProductMiddleware{usecase: usecase, log: log}
}

func (m *ProductMiddleware) Authorize(c *fiber.Ctx) error {
	// userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	endpoint := c.Route().Path
	switch endpoint {
	case "/api/products/:id", "/api/product-categories/:id":
		if role != "user" && role != "admin" && role != "super_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
		}
	case "/api/product-stocks/:id", "/api/warehouse-locations/:id":
		if role != "super_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Only super_admin can modify warehouse locations"})
		}
	default:
		if role != "admin" && role != "super_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: Only admin or super_admin can perform this action"})
		}
	}

	return c.Next()
}
