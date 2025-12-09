package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	"gorm.io/gorm"
	_ "github.com/jackc/pgx/v5/stdlib"

)

func EnsureDBExists(dsn string) {
	rootDSN, dbName := getRootDSNAndDBName(dsn)

	// Hubungkan ke 'postgres' database (database default)
	db, err := sql.Open("pgx", rootDSN)
	if err != nil {
		log.Fatal("Failed to open connection to root DB:", err)
	}
	defer db.Close()

	// Cek koneksi
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping root DB:", err)
	}

	// Cek apakah database target sudah ada
	var exists int
	checkQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName)
	err = db.QueryRow(checkQuery).Scan(&exists)

	if err == sql.ErrNoRows {
		log.Printf("Database '%s' does not exist. Creating...", dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Printf("Database '%s' created successfully", dbName)
	} else if err != nil {
		log.Printf("Warning: checking db existence failed: %v", err)
	} else {
		log.Printf("Database '%s' already exists.", dbName)
	}
}

// getRootDSNAndDBName memparsing DSN untuk mendapatkan nama DB target
// dan mengembalikan connection string baru yang mengarah ke database 'postgres'
func getRootDSNAndDBName(dsn string) (string, string) {
	// 1. Handle format URL (postgres://...)
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		if err != nil {
			log.Fatal("Invalid DSN URL:", err)
		}
		
		dbName := strings.TrimPrefix(u.Path, "/")
		if dbName == "" {
			dbName = "khalif_comment_db"
		}

		// Ubah path ke /postgres untuk koneksi root
		u.Path = "/postgres"
		return u.String(), dbName
	}

	// 2. Handle format Key-Value (host=... dbname=...)
	dbName := "khalif_comment_db" // default
	parts := strings.Split(dsn, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			dbName = strings.TrimPrefix(part, "dbname=")
		}
	}
	
	// Ganti dbname target dengan dbname=postgres
	rootDSN := strings.Replace(dsn, "dbname="+dbName, "dbname=postgres", 1)
	return rootDSN, dbName
}

func ResetSchema(db *gorm.DB) {
	queries := []string{
		"DROP SCHEMA public CASCADE;",
		"CREATE SCHEMA public;",
		"GRANT ALL ON SCHEMA public TO public;",
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Fatal(err)
		}
	}
}