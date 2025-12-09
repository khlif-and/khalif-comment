package main

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"khalif-comment/internal/config"
	"khalif-comment/pkg/database"

)

func ProvideDB(cfg *config.Config) *gorm.DB {
	// Memastikan database 'khalif_comment_db' ada sebelum connect
	database.EnsureDBExists(cfg.DBUrl)

	dbLogger := logger.Default.LogMode(logger.Error)

	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
		Logger:      dbLogger,
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// Tuning connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func ProvideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
}