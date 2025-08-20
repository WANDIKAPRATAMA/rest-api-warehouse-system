// TODO:
package controller

import (
	"auth-service/internal/dtos"
	"auth-service/internal/usecases"
	"auth-service/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthController interface {
	Signup(c *fiber.Ctx) error
	Signin(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	ChangeRole(c *fiber.Ctx) error
	Signout(c *fiber.Ctx) error
}

type authController struct {
	usecase  usecases.AuthUseCase
	log      *logrus.Logger
	validate *validator.Validate
}

func NewAuthController(usecase usecases.AuthUseCase, log *logrus.Logger, validate *validator.Validate) AuthController {
	return &authController{usecase: usecase, log: log, validate: validate}
}

func (c *authController) Signup(ctx *fiber.Ctx) error {
	var req dtos.SignupRequest
	allowedFields := utils.GenerateAllowedFields(dtos.SignupRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	user, err := c.usecase.Signup(ctx.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.StatusCreated, "User created successfully", fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	}, nil))
}

func (c *authController) Signin(ctx *fiber.Ctx) error {
	var req dtos.SigninRequest
	allowedFields := utils.GenerateAllowedFields(dtos.SigninRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}
	deviceID := ctx.Get("X-Device-ID")
	if deviceID == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse(
			fiber.StatusUnprocessableEntity,
			"Validation failed",
			[]utils.ErrorDetail{{Field: "X-Device-ID", Message: "Device ID required"}},
		))
	}
	accessToken, refreshToken, err := c.usecase.Signin(ctx.Context(), req.Email, req.Password, &deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Login successful", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil))
}

func (c *authController) ChangePassword(ctx *fiber.Ctx) error {
	var req dtos.ChangePasswordRequest
	allowedFields := utils.GenerateAllowedFields(dtos.ChangePasswordRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	if err := c.usecase.ChangePassword(ctx.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Password changed successfully", nil, nil))
}

func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	var req dtos.RefreshTokenRequest
	allowedFields := utils.GenerateAllowedFields(dtos.RefreshTokenRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	deviceID := ctx.Get("X-Device-ID")
	if deviceID == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorResponse(
			fiber.StatusUnprocessableEntity,
			"Validation failed",
			[]utils.ErrorDetail{{Field: "X-Device-ID", Message: "Device ID required"}},
		))
	}

	newAccessToken, newRefreshToken, err := c.usecase.RefreshToken(ctx.Context(), req.RefreshToken, deviceID)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(fiber.StatusUnauthorized, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Token refreshed successfully", fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken, // optional, bisa ganti token lama
	}, nil))
}

func (c *authController) ChangeRole(ctx *fiber.Ctx) error {
	var req dtos.ChangeRoleRequest
	allowedFields := utils.GenerateAllowedFields(dtos.ChangeRoleRequest{})
	if err := utils.BindAndValidateBody(ctx, &req, allowedFields, c.validate); err != nil {
		var errors []utils.ErrorDetail
		if validationErr := c.validate.Struct(req); validationErr != nil {
			for _, e := range validationErr.(validator.ValidationErrors) {
				errors = append(errors, utils.ErrorDetail{
					Field:   e.Field(),
					Message: e.Error(),
				})
			}
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error(), errors))
	}

	userID := ctx.Locals("userID").(uuid.UUID)
	if err := c.usecase.ChangeRole(ctx.Context(), userID, req.Role); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Role changed successfully", nil, nil))
}

func (c *authController) Signout(ctx *fiber.Ctx) error {
	tokenHash := ctx.Get("Authorization") // Simulasi, ganti dengan ekstrak dari header
	if err := c.usecase.Signout(ctx.Context(), tokenHash); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error(), nil))
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, "Signout successful", nil, nil))
}
