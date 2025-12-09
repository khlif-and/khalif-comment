package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"khalif-comment/internal/config"
	"khalif-comment/internal/domain"
	"khalif-comment/internal/handler"
	"khalif-comment/pkg/database"
	"khalif-comment/pkg/logger"

)

// @title           Khalif Comment API
// @version         1.0
// @description     API Service for User Comments
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@khalifstories.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8083
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
type App struct {
	DB             *gorm.DB
	RDB            *redis.Client
	CommentHandler *handler.CommentHandler
}

// NewApp hanya menerima CommentHandler
func NewApp(db *gorm.DB, rdb *redis.Client, ch *handler.CommentHandler) *App {
	return &App{
		DB:             db,
		RDB:            rdb,
		CommentHandler: ch,
	}
}

func main() {
	logger.Init()

	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()
	
	// InitializeApp dipanggil dari wire_gen.go (Dependency Injection)
	app, err := InitializeApp()
	if err != nil {
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}

	if *refreshFlag {
		database.ResetSchema(app.DB)
		logger.Info("Database reset successfully")
	}

	// AutoMigrate hanya untuk entitas Comment
	if err := app.DB.AutoMigrate(&domain.Comment{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// Menjalankan migrasi SQL tambahan (seperti trigger updated_at)
	database.RunMigrations(app.DB)

	// Seeder dihapus karena komentar tidak butuh data awal

	r := gin.New()
	r.Use(gin.Recovery())

	SetupRoutes(r, app, cfg)

	logger.Info("Server starting", zap.String("port", cfg.Port))
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Server start failed", zap.Error(err))
	}
}