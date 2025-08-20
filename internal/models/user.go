package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email         string         `gorm:"type:varchar(255);unique;not null"`
	Status        string         `gorm:"type:user_status;not null;default:'inactive'"`
	EmailVerified bool           `gorm:"column:email_verified;not null;default:false"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time      `gorm:"default:current_timestamp"`
	DeletedAt     gorm.DeletedAt `gorm:"index"` // Soft delete
}

type UserProfile struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null"`
	FullName     string         `gorm:"type:varchar(255)"`
	Phone        string         `gorm:"type:varchar(50)"`
	AvatarURL    string         `gorm:"type:text"`
	Address      string         `gorm:"type:text"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type UserSecurity struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null"`
	Password     string         `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type ApplicationRole struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID      `gorm:"column:source_user_id;type:uuid;not null"`
	Role         string         `gorm:"type:app_role;not null;default:'user'"`
	CreatedAt    time.Time      `gorm:"default:current_timestamp"`
	UpdatedAt    time.Time      `gorm:"default:current_timestamp"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type RefreshToken struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SourceUserID uuid.UUID `gorm:"column:source_user_id;type:uuid;not null;uniqueIndex:idx_user_device"`
	DeviceID     string    `gorm:"type:text;not null;uniqueIndex:idx_user_device"`
	TokenHash    string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
	ExpiresAt    time.Time `gorm:"not null"`
	LastUsedAt   time.Time
	RevokedAt    *time.Time     `gorm:"column:revoked_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
