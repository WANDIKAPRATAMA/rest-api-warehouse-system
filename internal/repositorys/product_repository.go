package repositorys

import (
	"auth-service/internal/dtos"
	"auth-service/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

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

	CreateWarehouseLocation(location *models.WarehouseLocation) error
	GetWarehouseLocationByID(id uuid.UUID) (*models.WarehouseLocation, error)
	UpdateWarehouseLocation(location *models.WarehouseLocation) error
	DeleteWarehouseLocation(id uuid.UUID) error

	CreateStockMovement(movement *models.StockMovement) error

	GetProductsList(req dtos.PaginationRequest) ([]models.Product, int64, error)
	GetWarehouseLocationsList(req dtos.PaginationRequest) ([]models.WarehouseLocation, int64, error)
	GetProductStocksList(req dtos.PaginationRequest) ([]models.ProductStock, int64, error)
	GetDashboardSummary() (*dtos.DashboardResponse, error)
	GetProductCategoriesList(req dtos.PaginationRequest) ([]models.ProductCategory, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

// CreateWarehouseLocation
func (r *productRepository) CreateWarehouseLocation(location *models.WarehouseLocation) error {
	return r.db.Create(location).Error
}

// GetWarehouseLocationByID
func (r *productRepository) GetWarehouseLocationByID(id uuid.UUID) (*models.WarehouseLocation, error) {
	var location models.WarehouseLocation
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&location).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

// UpdateWarehouseLocation
func (r *productRepository) UpdateWarehouseLocation(location *models.WarehouseLocation) error {
	return r.db.Save(location).Error
}

// DeleteWarehouseLocation
func (r *productRepository) DeleteWarehouseLocation(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.WarehouseLocation{}).Error
}

// Todo: Warhouse Implemetation

func (r *productRepository) GetProductCategoriesList(req dtos.PaginationRequest) ([]models.ProductCategory, int64, error) {
	var categories []models.ProductCategory
	var total int64

	query := r.db.Model(&models.ProductCategory{}).Where("deleted_at IS NULL")
	if req.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
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

	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
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
func (r *productRepository) GetDashboardSummary() (*dtos.DashboardResponse, error) {
	summary := &dtos.DashboardResponse{}

	// Total stock: SUM(quantity)
	r.db.Model(&models.ProductStock{}).Select("SUM(quantity)").Where("deleted_at IS NULL").Scan(&summary.TotalStock)

	// Number of products: COUNT(distinct products)
	r.db.Model(&models.Product{}).Where("deleted_at IS NULL").Count(&summary.NumberOfProducts)

	// Low-stock items detail (join product, warehouse, user, profile)
	var lowStockItems []struct {
		ProductID      uuid.UUID
		ProductName    string
		WarehouseID    uuid.UUID
		WarehouseName  string
		Quantity       int
		Status         string
		UpdatedByEmail string
		UpdatedByName  string
		UpdatedAt      time.Time
	}
	r.db.Table("product_stocks ps").
		Select("ps.id as product_id, p.name as product_name, ps.warehouse_location_id as warehouse_id, wl.name as warehouse_name, ps.quantity, ps.status, u.email as updated_by_email, up.full_name as updated_by_name, ps.updated_at").
		Joins("JOIN products p ON p.id = ps.source_product_id").
		Joins("JOIN warehouse_locations wl ON wl.id = ps.warehouse_location_id").
		Joins("JOIN users u ON u.id = ps.updated_by").
		Joins("JOIN user_profiles up ON up.source_user_id = u.id").
		Where("ps.status = 'low-stock' AND ps.deleted_at IS NULL").
		Limit(10). // Batasi untuk performa, misal top 10
		Scan(&lowStockItems)
	for _, item := range lowStockItems {
		summary.LowStockItems = append(summary.LowStockItems, dtos.LowStockDetail{
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			WarehouseID:    item.WarehouseID,
			WarehouseName:  item.WarehouseName,
			Quantity:       item.Quantity,
			Status:         item.Status,
			UpdatedByEmail: item.UpdatedByEmail,
			UpdatedByName:  item.UpdatedByName,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	// Out-of-stock items detail (serupa dengan low-stock)
	var outOfStockItems []struct {
		ProductID      uuid.UUID
		ProductName    string
		WarehouseID    uuid.UUID
		WarehouseName  string
		Quantity       int
		Status         string
		UpdatedByEmail string
		UpdatedByName  string
		UpdatedAt      time.Time
	}
	r.db.Table("product_stocks ps").
		Select("ps.id as product_id, p.name as product_name, ps.warehouse_location_id as warehouse_id, wl.name as warehouse_name, ps.quantity, ps.status, u.email as updated_by_email, up.full_name as updated_by_name, ps.updated_at").
		Joins("JOIN products p ON p.id = ps.source_product_id").
		Joins("JOIN warehouse_locations wl ON wl.id = ps.warehouse_location_id").
		Joins("JOIN users u ON u.id = ps.updated_by").
		Joins("JOIN user_profiles up ON up.source_user_id = u.id").
		Where("ps.status = 'out-of-stock' AND ps.deleted_at IS NULL").
		Limit(10).
		Scan(&outOfStockItems)
	for _, item := range outOfStockItems {
		summary.OutOfStockItems = append(summary.OutOfStockItems, dtos.OutOfStockDetail{
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			WarehouseID:    item.WarehouseID,
			WarehouseName:  item.WarehouseName,
			Quantity:       item.Quantity,
			Status:         item.Status,
			UpdatedByEmail: item.UpdatedByEmail,
			UpdatedByName:  item.UpdatedByName,
			UpdatedAt:      item.UpdatedAt,
		})
	}

	var recentAdditions []struct {
		ProductID      uuid.UUID
		ProductName    string
		CreatedByEmail string
		CreatedByName  string
		CreatedAt      time.Time
	}
	r.db.Table("products p").
		Select("p.id as product_id, p.name as product_name, u.email as created_by_email, up.full_name as created_by_name, p.created_at").
		Joins("JOIN users u ON u.id = p.created_by").
		Joins("JOIN user_profiles up ON up.source_user_id = u.id").
		Where("p.deleted_at IS NULL").
		Order("p.created_at DESC").
		Limit(5).
		Scan(&recentAdditions)
	for _, item := range recentAdditions {
		summary.RecentAdditions = append(summary.RecentAdditions, dtos.RecentAddition{
			ProductID:      item.ProductID,
			ProductName:    item.ProductName,
			CreatedByEmail: item.CreatedByEmail,
			CreatedByName:  item.CreatedByName,
			CreatedAt:      item.CreatedAt,
		})
	}

	return summary, nil
}
