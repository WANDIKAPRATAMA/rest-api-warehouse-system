package repositorys

import (
	"auth-service/internal/dtos"
	"auth-service/internal/models"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uuid.UUID) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uuid.UUID) error
	CreateProductCategory(category *models.ProductCategory) error
	GetProductCategoryByID(id uuid.UUID) (*models.ProductCategory, error)
	UpdateProductCategory(category *models.ProductCategory) error
	DeleteProductCategory(id uuid.UUID) error
	CreateProductStock(stock *models.ProductStock) error
	GetProductStockByID(id uuid.UUID) (*models.ProductStock, error)
	UpdateProductStock(stock *models.ProductStock) error
	DeleteProductStock(id uuid.UUID) error
	CreateStockMovement(movement *models.StockMovement) error

	GetProductsList(req dtos.PaginationRequest) ([]models.Product, int64, error)
	GetWarehouseLocationsList(req dtos.PaginationRequest) ([]models.WarehouseLocation, int64, error)
	GetProductStocksList(req dtos.PaginationRequest) ([]models.ProductStock, int64, error)
	GetDashboardSummary() (dtos.DashboardResponse, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateStockMovement(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) DeleteProduct(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.Product{}).Error
}

var ErrCategoryNameAlreadyExists = errors.New("category name already exists")

func (r *productRepository) CreateProductCategory(category *models.ProductCategory) error {
	err := r.db.Create(category).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "name") {
			return ErrCategoryNameAlreadyExists
		}
	}
	return err
}

func (r *productRepository) GetProductCategoryByID(id uuid.UUID) (*models.ProductCategory, error) {
	var category models.ProductCategory
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *productRepository) UpdateProductCategory(category *models.ProductCategory) error {
	return r.db.Save(category).Error
}

func (r *productRepository) DeleteProductCategory(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.ProductCategory{}).Error
}

func (r *productRepository) CreateProductStock(stock *models.ProductStock) error {
	return r.db.Create(stock).Error
}

func (r *productRepository) GetProductStockByID(id uuid.UUID) (*models.ProductStock, error) {
	var stock models.ProductStock
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *productRepository) UpdateProductStock(stock *models.ProductStock) error {
	return r.db.Save(stock).Error
}

func (r *productRepository) DeleteProductStock(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.ProductStock{}).Error
}

// GetProductsList dengan join, filter, search, pagination
func (r *productRepository) GetProductsList(req dtos.PaginationRequest) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("deleted_at IS NULL")
	if req.CategoryID != uuid.Nil {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	// Hitung total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination dan sorting
	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit).Offset(offset)
	if req.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, req.Order))
	}

	// Join/Preload category
	query = query.Preload("Category")

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// GetWarehouseLocationsList
func (r *productRepository) GetWarehouseLocationsList(req dtos.PaginationRequest) ([]models.WarehouseLocation, int64, error) {
	var locations []models.WarehouseLocation
	var total int64

	query := r.db.Model(&models.WarehouseLocation{}).Where("deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit).Offset(offset)
	if req.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, req.Order))
	}

	if err := query.Find(&locations).Error; err != nil {
		return nil, 0, err
	}
	return locations, total, nil
}

// GetProductStocksList dengan join product dan warehouse
func (r *productRepository) GetProductStocksList(req dtos.PaginationRequest) ([]models.ProductStock, int64, error) {
	var stocks []models.ProductStock
	var total int64

	query := r.db.Model(&models.ProductStock{}).Where("deleted_at IS NULL")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Search != "" {
		query = query.Joins("JOIN products ON products.id = product_stocks.source_product_id").
			Where("products.name ILIKE ?", "%"+req.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.Limit
	query = query.Limit(req.Limit).Offset(offset)
	if req.SortBy != "" {
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, req.Order))
	}

	// Preload relasi
	query = query.Preload("Product").Preload("WarehouseLocation")

	if err := query.Find(&stocks).Error; err != nil {
		return nil, 0, err
	}
	return stocks, total, nil
}

// GetDashboardSummary
func (r *productRepository) GetDashboardSummary() (dtos.DashboardResponse, error) {
	var summary dtos.DashboardResponse

	// Total stock: SUM(quantity)
	r.db.Model(&models.ProductStock{}).Select("SUM(quantity)").Where("deleted_at IS NULL").Scan(&summary.TotalStock)

	// Number of products: COUNT(distinct products)
	r.db.Model(&models.Product{}).Where("deleted_at IS NULL").Count(&summary.NumberOfProducts)

	// Low-stock items
	r.db.Model(&models.ProductStock{}).Where("status = 'low-stock' AND deleted_at IS NULL").Count(&summary.LowStockItems)

	// Out-of-stock items
	r.db.Model(&models.ProductStock{}).Where("status = 'out-of-stock' AND deleted_at IS NULL").Count(&summary.OutOfStockItems)

	return summary, nil
}
