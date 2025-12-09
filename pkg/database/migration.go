package database

import (
	"embed"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"khalif-comment/pkg/logger" // Ubah import path sesuai nama module

)

//go:embed schema/*.sql
var schemaFS embed.FS

func RunMigrations(db *gorm.DB) {
	// Pastikan file 'schema/001_init_schema.sql' ada di folder pkg/database/schema/
	content, err := schemaFS.ReadFile("schema/001_init_schema.sql")
	if err != nil {
		logger.Fatal("Failed to read migration file from embed", zap.Error(err))
	}

	blocks := strings.Split(string(content), "--SEPARATOR--")

	for _, block := range blocks {
		trimmedBlock := strings.TrimSpace(block)
		if trimmedBlock == "" {
			continue
		}

		if err := db.Exec(trimmedBlock).Error; err != nil {
			logger.Fatal("Failed to execute migration block", 
				zap.String("query_snippet", trimmedBlock[:min(len(trimmedBlock), 50)]), 
				zap.Error(err),
			)
		}
	}

	logger.Info("Database migration executed successfully")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}