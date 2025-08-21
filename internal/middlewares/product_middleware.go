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
	role := c.Locals("role").(string)
	endpoint := c.Route().Path
	method := c.Method() // GET, POST, PUT, DELETE

	// Rule khusus berdasarkan method
	switch method {
	case fiber.MethodGet:
		// Alow users role to GET (list atau detail)
		return c.Next()

	case fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete:
		if role == "user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: user role is only allowed to view data",
			})
		}
	}

	switch endpoint {
	case "/api/product-stocks/:id", "/api/warehouse-locations/:id":
		if method != fiber.MethodGet && role != "super_admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden: Only super_admin can modify warehouse locations or product stocks",
			})
		}
	}

	return c.Next()
}
