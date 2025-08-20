package routes

import (
	"auth-service/internal/controllers"
	middleware "auth-service/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

type ProductRouteConfig struct {
	App               *fiber.App
	ProductController controllers.ProductController
	ProductMiddleware *middleware.ProductMiddleware
	AuthMiddleware    *middleware.AuthMiddleware
}

func (r *ProductRouteConfig) Setup() {
	api := r.App.Group("/api")

	products := api.Group("/products", r.AuthMiddleware.Authenticate)
	products.Post("/", r.ProductMiddleware.Authorize, r.ProductController.CreateProduct)
	products.Get("/:id", r.ProductMiddleware.Authorize, r.ProductController.GetProductByID)
	products.Put("/:id", r.ProductMiddleware.Authorize, r.ProductController.UpdateProduct)
	products.Delete("/:id", r.ProductMiddleware.Authorize, r.ProductController.DeleteProduct)
	products.Get("/", r.ProductMiddleware.Authorize, r.ProductController.GetProductsList)

	categories := api.Group("/product-categories", r.AuthMiddleware.Authenticate)
	categories.Post("/", r.ProductMiddleware.Authorize, r.ProductController.CreateProductCategory)
	categories.Get("/:id", r.ProductMiddleware.Authorize, r.ProductController.GetProductCategoryByID)
	categories.Put("/:id", r.ProductMiddleware.Authorize, r.ProductController.UpdateProductCategory)
	categories.Delete("/:id", r.ProductMiddleware.Authorize, r.ProductController.DeleteProductCategory)

	stocks := api.Group("/product-stocks", r.AuthMiddleware.Authenticate)
	stocks.Post("/", r.ProductMiddleware.Authorize, r.ProductController.CreateProductStock)
	stocks.Get("/:id", r.ProductMiddleware.Authorize, r.ProductController.GetProductStockByID)
	stocks.Put("/:id", r.ProductMiddleware.Authorize, r.ProductController.UpdateProductStock)
	stocks.Delete("/:id", r.ProductMiddleware.Authorize, r.ProductController.DeleteProductStock)
	stocks.Get("/", r.ProductMiddleware.Authorize, r.ProductController.GetProductStocksList)

	warehouse := api.Group("/warehouse-locations", r.AuthMiddleware.Authenticate)
	warehouse.Get("/", r.ProductMiddleware.Authorize, r.ProductController.GetWarehouseLocationsList)

	dashboard := api.Group("/dashboard", r.AuthMiddleware.Authenticate)
	dashboard.Get("/", r.ProductMiddleware.Authorize, r.ProductController.GetDashboardSummary)
}
