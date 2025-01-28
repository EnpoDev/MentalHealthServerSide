package main

import (
	"log"
	"mental-health-companion/internal/database"
	"mental-health-companion/internal/handlers"
	"mental-health-companion/internal/middleware"
	"mental-health-companion/internal/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Veritabanı bağlantısını başlat
	database.InitDB()

	// Tabloları otomatik oluştur
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Gin router'ı oluştur
	r := gin.Default()

	// CORS ayarları
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rate limiting middleware'i ekle
	r.Use(rateLimiter())

	// Public rotalar
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected rotalar
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handlers.Me)
	}

	// Sunucuyu başlat
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func rateLimiter() gin.HandlerFunc {
	// TODO: Rate limiting implementasyonu eklenecek
	return func(c *gin.Context) {
		c.Next()
	}
} 