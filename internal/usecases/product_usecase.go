package usecases

import (
	"context"
	"fmt"
	"time"

	"auth-service/internal/dtos"
	"auth-service/internal/models"
	"auth-service/internal/repositorys"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductUseCase interface {
	CreateProduct(ctx context.Context, req dtos.CreateProductRequest, userID uuid.UUID) (*dtos.ProductResponse, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*dtos.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req dtos.UpdateProductRequest, userID uuid.UUID) (*dtos.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CreateProductCategory(ctx context.Context, req dtos.CreateProductCategoryRequest, userID uuid.UUID) (*dtos.ProductCategoryResponse, error)
	GetProductCategoryByID(ctx context.Context, id uuid.UUID) (*dtos.ProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, id uuid.UUID, req dtos.UpdateProductCategoryRequest, userID uuid.UUID) (*dtos.ProductCategoryResponse, error)
	DeleteProductCategory(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	CreateProductStock(ctx context.Context, req dtos.CreateProductStockRequest, userID uuid.UUID) (*dtos.ProductStockResponse, error)
	GetProductStockByID(ctx context.Context, id uuid.UUID) (*dtos.ProductStockResponse, error)
	UpdateProductStock(ctx context.Context, id uuid.UUID, req dtos.UpdateProductStockRequest, userID uuid.UUID) (*dtos.ProductStockResponse, error)
	DeleteProductStock(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	TrackStockMovement(ctx context.Context, productID uuid.UUID, movementType string, quantity int, userID uuid.UUID) error

	GetProductsList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.ProductListResponse, dtos.Pagination, error)
	GetWarehouseLocationsList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.WarehouseLocationListResponse, dtos.Pagination, error)
	GetProductStocksList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.ProductStockListResponse, dtos.Pagination, error)
	GetDashboardSummary(ctx context.Context) (*dtos.DashboardResponse, error)
}

type productUseCase struct {
	repo     repositorys.ProductRepository
	validate *validator.Validate
	log      *logrus.Logger
}

func NewProductUseCase(repo repositorys.ProductRepository, log *logrus.Logger, validate *validator.Validate) ProductUseCase {
	return &productUseCase{repo: repo, log: log, validate: validate}
}

func (u *productUseCase) CreateProduct(ctx context.Context, req dtos.CreateProductRequest, userID uuid.UUID) (*dtos.ProductResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	product := &models.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		SKU:         req.SKU,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		CreatedBy:   userID,
	}
	if err := u.repo.CreateProduct(product); err != nil {
		return nil, err
	}
	u.log.Info(fmt.Sprintf("Product %s created", product.Name))
	return &dtos.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		SKU:         product.SKU,
		CategoryID:  product.CategoryID,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) GetProductByID(ctx context.Context, id uuid.UUID) (*dtos.ProductResponse, error) {
	product, err := u.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return &dtos.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		SKU:         product.SKU,
		CategoryID:  product.CategoryID,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) UpdateProduct(ctx context.Context, id uuid.UUID, req dtos.UpdateProductRequest, userID uuid.UUID) (*dtos.ProductResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	product, err := u.repo.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.CategoryID = req.CategoryID
	product.Description = req.Description
	product.UpdatedAt = time.Now()
	if err := u.repo.UpdateProduct(product); err != nil {
		return nil, err
	}

	return &dtos.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		SKU:         product.SKU,
		CategoryID:  product.CategoryID,
		Description: product.Description,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) DeleteProduct(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := u.repo.GetProductByID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteProduct(id)
}

// Implementasi serupa untuk ProductCategory
func (u *productUseCase) CreateProductCategory(ctx context.Context, req dtos.CreateProductCategoryRequest, userID uuid.UUID) (*dtos.ProductCategoryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	category := &models.ProductCategory{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
	}
	if err := u.repo.CreateProductCategory(category); err != nil {
		return nil, err
	}

	return &dtos.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) GetProductCategoryByID(ctx context.Context, id uuid.UUID) (*dtos.ProductCategoryResponse, error) {
	category, err := u.repo.GetProductCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return &dtos.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) UpdateProductCategory(ctx context.Context, id uuid.UUID, req dtos.UpdateProductCategoryRequest, userID uuid.UUID) (*dtos.ProductCategoryResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	category, err := u.repo.GetProductCategoryByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	category.UpdatedAt = time.Now()
	if err := u.repo.UpdateProductCategory(category); err != nil {
		return nil, err
	}

	return &dtos.ProductCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) DeleteProductCategory(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := u.repo.GetProductCategoryByID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteProductCategory(id)
}

// Implementasi untuk ProductStock
func (u *productUseCase) CreateProductStock(ctx context.Context, req dtos.CreateProductStockRequest, userID uuid.UUID) (*dtos.ProductStockResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	stock := &models.ProductStock{
		ID:                  uuid.New(),
		SourceProductID:     req.ProductID,
		WarehouseLocationID: req.WarehouseLocationID,
		Quantity:            req.Quantity,
		UpdatedBy:           userID,
	}
	if err := u.repo.CreateProductStock(stock); err != nil {
		return nil, err
	}

	return &dtos.ProductStockResponse{
		ID:                  stock.ID,
		ProductID:           stock.SourceProductID,
		WarehouseLocationID: stock.WarehouseLocationID,
		Quantity:            stock.Quantity,
		Status:              stock.Status,
		UpdatedAt:           stock.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) GetProductStockByID(ctx context.Context, id uuid.UUID) (*dtos.ProductStockResponse, error) {
	stock, err := u.repo.GetProductStockByID(id)
	if err != nil {
		return nil, err
	}
	return &dtos.ProductStockResponse{
		ID:                  stock.ID,
		ProductID:           stock.SourceProductID,
		WarehouseLocationID: stock.WarehouseLocationID,
		Quantity:            stock.Quantity,
		Status:              stock.Status,
		UpdatedAt:           stock.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) UpdateProductStock(ctx context.Context, id uuid.UUID, req dtos.UpdateProductStockRequest, userID uuid.UUID) (*dtos.ProductStockResponse, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, err
	}

	stock, err := u.repo.GetProductStockByID(id)
	if err != nil {
		return nil, err
	}

	stock.Quantity = req.Quantity
	stock.UpdatedAt = time.Now()
	stock.UpdatedBy = userID
	if err := u.repo.UpdateProductStock(stock); err != nil {
		return nil, err
	}

	return &dtos.ProductStockResponse{
		ID:                  stock.ID,
		ProductID:           stock.SourceProductID,
		WarehouseLocationID: stock.WarehouseLocationID,
		Quantity:            stock.Quantity,
		Status:              stock.Status,
		UpdatedAt:           stock.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *productUseCase) DeleteProductStock(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := u.repo.GetProductStockByID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteProductStock(id)
}

func (u *productUseCase) TrackStockMovement(ctx context.Context, productID uuid.UUID, movementType string, quantity int, userID uuid.UUID) error {
	// Validasi movementType
	validMovements := map[string]bool{"inbound": true, "outbound": true}
	if !validMovements[movementType] {
		return fmt.Errorf("invalid movement type: %s", movementType)
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Ambil stok terkait
	stock, err := u.repo.GetProductStockByID(productID) // Asumsi ada logika untuk mencari stok berdasarkan productID
	if err != nil {
		return fmt.Errorf("stock not found for product ID: %v", err)
	}

	// Hitung stok baru
	newQuantity := stock.Quantity
	if movementType == "inbound" {
		newQuantity += quantity
	} else if movementType == "outbound" {
		newQuantity -= quantity
		if newQuantity < 0 {
			return fmt.Errorf("insufficient stock for outbound movement")
		}
	}

	// Perbarui stok
	stock.Quantity = newQuantity
	stock.Status = u.determineStockStatus(newQuantity)
	stock.UpdatedAt = time.Now()
	stock.UpdatedBy = userID
	if err := u.repo.UpdateProductStock(stock); err != nil {
		return err
	}

	// Catat pergerakan stok
	movement := &models.StockMovement{
		ID:              uuid.New(),
		SourceProductID: productID,
		MovementType:    movementType,
		Quantity:        quantity,
		ReferenceNote:   "Automatic stock update",
		CreatedBy:       userID,
		CreatedAt:       time.Now(),
	}
	return u.repo.CreateStockMovement(movement)
}

func (u *productUseCase) determineStockStatus(quantity int) string {
	switch {
	case quantity <= 0:
		return "out-of-stock"
	case quantity < 10:
		return "low-stock"
	default:
		return "available"
	}
}

func (u *productUseCase) GetProductsList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.ProductListResponse, dtos.Pagination, error) {
	if err := u.validate.Struct(req); err != nil {
		return nil, dtos.Pagination{}, err
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Order == "" {
		req.Order = "asc"
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}

	products, total, err := u.repo.GetProductsList(req)
	if err != nil {
		return nil, dtos.Pagination{}, err
	}

	var list []dtos.ProductListResponse
	for _, p := range products {
		list = append(list, dtos.ProductListResponse{
			ID:           p.ID,
			Name:         p.Name,
			SKU:          p.SKU,
			CategoryID:   p.CategoryID,
			CategoryName: p.Category.Name, // Dari preload
			Description:  p.Description,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	hasNextPage := req.Page < totalPages
	nextPage := req.Page + 1
	if !hasNextPage {
		nextPage = 0
	}

	pagination := dtos.Pagination{
		HasNextPage: hasNextPage,
		NextPage:    &nextPage,
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		TotalItems:  int(total),
	}

	return list, pagination, nil
}

func (u *productUseCase) GetWarehouseLocationsList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.WarehouseLocationListResponse, dtos.Pagination, error) {
	// Logika serupa dengan GetProductsList, adaptasi untuk WarehouseLocation
	locations, total, err := u.repo.GetWarehouseLocationsList(req)
	if err != nil {
		return nil, dtos.Pagination{}, err
	}

	var list []dtos.WarehouseLocationListResponse
	for _, l := range locations {
		list = append(list, dtos.WarehouseLocationListResponse{
			ID:          l.ID,
			Name:        l.Name,
			Description: l.Description,
			CreatedAt:   l.CreatedAt,
		})
	}

	// Hitung pagination seperti di atas
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	hasNextPage := req.Page < totalPages
	nextPage := req.Page + 1
	if !hasNextPage {
		nextPage = 0
	}

	pagination := dtos.Pagination{
		HasNextPage: hasNextPage,
		NextPage:    &nextPage,
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		TotalItems:  int(total),
	}

	return list, pagination, nil
}

func (u *productUseCase) GetProductStocksList(ctx context.Context, req dtos.PaginationRequest) ([]dtos.ProductStockListResponse, dtos.Pagination, error) {

	stocks, total, err := u.repo.GetProductStocksList(req)
	if err != nil {
		return nil, dtos.Pagination{}, err
	}

	var list []dtos.ProductStockListResponse
	for _, s := range stocks {
		list = append(list, dtos.ProductStockListResponse{
			ID:                  s.ID,
			ProductID:           s.SourceProductID,
			ProductName:         s.Product.Name, // Dari preload
			WarehouseLocationID: s.WarehouseLocationID,
			WarehouseName:       s.WarehouseLocation.Name, // Dari preload
			Quantity:            s.Quantity,
			Status:              s.Status,
			UpdatedAt:           s.UpdatedAt,
		})
	}

	// Hitung pagination seperti di atas
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	hasNextPage := req.Page < totalPages
	nextPage := req.Page + 1
	if !hasNextPage {
		nextPage = 0
	}

	pagination := dtos.Pagination{
		HasNextPage: hasNextPage,
		NextPage:    &nextPage,
		CurrentPage: req.Page,
		TotalPages:  totalPages,
		TotalItems:  int(total),
	}

	return list, pagination, nil
}

func (u *productUseCase) GetDashboardSummary(ctx context.Context) (*dtos.DashboardResponse, error) {
	summary, err := u.repo.GetDashboardSummary()
	if err != nil {
		return nil, err
	}
	return &summary, nil
}
