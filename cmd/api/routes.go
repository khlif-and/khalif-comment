package main

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "khalif-comment/docs" // Pastikan folder docs di-generate oleh swag init
	"khalif-comment/internal/config"
	"khalif-comment/pkg/middleware"

)

func SetupRoutes(r *gin.Engine, app *App, cfg *config.Config) {
	r.Use(middleware.Logger())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Rate Limiter (Menggunakan Redis dari App)
	limiter := middleware.RateLimitConfig{Limit: 300, Window: time.Minute}
	r.Use(middleware.RateLimit(app.RDB, limiter))
	
	auth := middleware.AuthMiddleware(cfg.JWTSecret)

	// --- Public Routes ---
	// User bisa membaca komentar tanpa harus login
	r.GET("/api/comments", app.CommentHandler.GetByStory)

	// --- Protected Routes (User Login) ---
	protected := r.Group("/api")
	protected.Use(auth)
	{
		// Create, Update, Delete butuh data UserID dari token
		protected.POST("/comments", app.CommentHandler.Create)
		protected.PUT("/comments/:id", app.CommentHandler.Update)
		protected.DELETE("/comments/:id", app.CommentHandler.Delete)
	}
}