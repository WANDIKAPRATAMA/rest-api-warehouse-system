package configs

import (
	"auth-service/internal/models"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, log *logrus.Logger) *gorm.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	sslmode := viper.GetString("database.sslmode")
	timezone := viper.GetString("database.timezone")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	// Format DSN PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		host, username, password, database, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 5,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Connection pool setup
	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection: %v", err)
	}
	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))
	enumTypes := map[string][]string{
		"user_status": {
			"active",
			"inactive",
			"banned",
		},
		"app_role": {
			"super_admin",
			"admin",
			"user",
		},
		"stock_status": {
			"available",
			"low-stock",
			"out-of-stock",
		},
		"movement_type": {
			"inbound",
			"outbound",
		},
	}
	for typeName, values := range enumTypes {
		var exists bool
		err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?)", typeName).Scan(&exists).Error
		if err != nil {
			log.Printf("gagal memeriksa tipe %s: %v", typeName, err)
		}
		if !exists {
			valuesStr := "'" + strings.Join(values, "', '") + "'"
			err = db.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", typeName, valuesStr)).Error
			if err != nil {
				log.Printf("Gagal membuat tipe %s: %v", typeName, err)
			}
		}
	}

	// Auto migrate models (urutan penting: parent dulu)
	err = db.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.UserSecurity{},
		&models.ApplicationRole{},
		&models.RefreshToken{},
		&models.ProductCategory{},
		&models.Product{},
		&models.WarehouseLocation{},
		&models.ProductStock{},
		&models.StockMovement{},
	)
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}
	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (l *logrusWriter) Printf(message string, args ...interface{}) {
	l.Logger.Tracef(message, args...)
}
