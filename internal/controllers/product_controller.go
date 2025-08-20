package controllers

import (
	"auth-service/internal/dtos"
	"auth-service/internal/usecases"
	"auth-service/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductController interface {
	CreateProduct(c *fiber.Ctx) error
	GetProductByID(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	CreateProductCategory(c *fiber.Ctx) error
	GetProductCategoryByID(c *fiber.Ctx) error
	UpdateProductCategory(c *fiber.Ctx) error
	DeleteProductCategory(c *fiber.Ctx) error
	CreateProductStock(c *fiber.Ctx) error
	GetProductStockByID(c *fiber.Ctx) error
	UpdateProductStock(c *fiber.Ctx) error
	DeleteProductStock(c *fiber.Ctx) error

	GetProductsList(ctx *fiber.Ctx) error
	GetWarehouseLocationsList(ctx *fiber.Ctx) error
	GetProductStocksList(ctx *fiber.Ctx) error
	GetDashboardSummary(ctx *fiber.Ctx) error
}

type productController struct {
	usecase  usecases.ProductUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewProductController(usecase usecases.ProductUseCase, log *logrus.Logger, validate *validator.Validate) ProductController {
	return &productController{usecase: usecase, log: log, validate: validate}
}

func (c *productController) CreateProductCategory(ctx *fiber.Ctx) error {
	var req dtos.CreateProductCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	category, err := c.usecase.CreateProductCategory(ctx.Context(), req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusCreated,
		Message:    "Product category created successfully",
		Payload:    category,
	})
}

func (c *productController) GetProductCategoryByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	category, err := c.usecase.GetProductCategoryByID(ctx.Context(), categoryID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusNotFound,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product category retrieved successfully",
		Payload:    category,
	})
}

func (c *productController) UpdateProductCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	var req dtos.UpdateProductCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	category, err := c.usecase.UpdateProductCategory(ctx.Context(), categoryID, req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product category updated successfully",
		Payload:    category,
	})
}

func (c *productController) DeleteProductCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	if err := c.usecase.DeleteProductCategory(ctx.Context(), categoryID, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product category deleted successfully",
		Payload:    nil,
	})
}

func (c *productController) CreateProductStock(ctx *fiber.Ctx) error {
	var req dtos.CreateProductStockRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	stock, err := c.usecase.CreateProductStock(ctx.Context(), req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusCreated,
		Message:    "Product stock created successfully",
		Payload:    stock,
	})
}

func (c *productController) GetProductStockByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	stockID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	stock, err := c.usecase.GetProductStockByID(ctx.Context(), stockID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusNotFound,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product stock retrieved successfully",
		Payload:    stock,
	})
}

func (c *productController) UpdateProductStock(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	stockID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	var req dtos.UpdateProductStockRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	stock, err := c.usecase.UpdateProductStock(ctx.Context(), stockID, req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product stock updated successfully",
		Payload:    stock,
	})
}

func (c *productController) DeleteProductStock(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	stockID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	if err := c.usecase.DeleteProductStock(ctx.Context(), stockID, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product stock deleted successfully",
		Payload:    nil,
	})
}

func (c *productController) CreateProduct(ctx *fiber.Ctx) error {
	var req dtos.CreateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	product, err := c.usecase.CreateProduct(ctx.Context(), req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusCreated,
		Message:    "Product created successfully",
		Payload:    product,
	})
}

func (c *productController) GetProductByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	product, err := c.usecase.GetProductByID(ctx.Context(), productID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusNotFound,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product retrieved successfully",
		Payload:    product,
	})
}

func (c *productController) UpdateProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	var req dtos.UpdateProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	product, err := c.usecase.UpdateProduct(ctx.Context(), productID, req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product updated successfully",
		Payload:    product,
	})
}

func (c *productController) DeleteProduct(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid ID format",
			Payload:    nil,
		})
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	if err := c.usecase.DeleteProduct(ctx.Context(), productID, userID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dtos.RefreshTokenResponse{
			Status:     "error",
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
			Payload:    nil,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(dtos.RefreshTokenResponse{
		Status:     "success",
		StatusCode: fiber.StatusOK,
		Message:    "Product deleted successfully",
		Payload:    nil,
	})
}

// Implementasi serupa untuk ProductCategory dan ProductStock

func (c *productController) GetProductsList(ctx *fiber.Ctx) error {
	var req dtos.PaginationRequest
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}

	list, pagination, err := c.usecase.GetProductsList(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Products list retrieved", list, pagination))
}

func (c *productController) GetWarehouseLocationsList(ctx *fiber.Ctx) error {
	var req dtos.PaginationRequest
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}

	list, pagination, err := c.usecase.GetWarehouseLocationsList(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Warehouse locations list retrieved", list, pagination))
}

func (c *productController) GetProductStocksList(ctx *fiber.Ctx) error {
	var req dtos.PaginationRequest
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}
	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), nil))
	}

	list, pagination, err := c.usecase.GetProductStocksList(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Product stocks list retrieved", list, pagination))
}

func (c *productController) GetDashboardSummary(ctx *fiber.Ctx) error {
	summary, err := c.usecase.GetDashboardSummary(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Dashboard summary retrieved", summary, nil))
}
