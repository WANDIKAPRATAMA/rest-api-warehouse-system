package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"type:varchar(100);unique;not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"default:current_timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"default:current_timestamp"`
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"not null"`
	SKU         string    `gorm:"unique;not null"`
	CategoryID  uuid.UUID `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CreatedBy   uuid.UUID      `gorm:"column:created_by;type:uuid"`

	// Relasi ke category
	Category ProductCategory `gorm:"foreignKey:CategoryID;references:ID"`
}

type WarehouseLocation struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"type:varchar(100);not null"`
	Description string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"default:current_timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"default:current_timestamp"`
}

type ProductStock struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceProductID     uuid.UUID      `gorm:"column:source_product_id;type:uuid;not null"`
	WarehouseLocationID uuid.UUID      `gorm:"column:warehouse_location_id;type:uuid;not null"`
	Quantity            int            `gorm:"not null;default:0"`
	Status              string         `gorm:"type:stock_status;default:'available'"`
	UpdatedBy           uuid.UUID      `gorm:"column:updated_by;type:uuid"`
	UpdatedAt           time.Time      `gorm:"default:current_timestamp"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	CreatedAt           time.Time      `gorm:"default:current_timestamp"`

	Product           Product           `gorm:"foreignKey:SourceProductID;references:ID"`
	WarehouseLocation WarehouseLocation `gorm:"foreignKey:WarehouseLocationID;references:ID"`
}

type StockMovement struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceProductID uuid.UUID      `gorm:"column:source_product_id;type:uuid;not null"`
	MovementType    string         `gorm:"type:movement_type;not null"`
	Quantity        int            `gorm:"not null"`
	ReferenceNote   string         `gorm:"type:text"`
	CreatedBy       uuid.UUID      `gorm:"column:created_by;type:uuid"`
	CreatedAt       time.Time      `gorm:"default:current_timestamp"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}
