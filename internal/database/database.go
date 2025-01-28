package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	config := &gorm.Config{
		Logger: newLogger,
		PrepareStmt: true, // Statement caching için
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Veritabanına bağlanılamadı:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("SQL DB alınamadı:", err)
	}

	// Bağlantı havuzu ayarları
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
} 