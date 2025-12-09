package config

import (
	"log"
	"os"

	"github.com/spf13/viper"

)

type Config struct {
	DBUrl     string `mapstructure:"DATABASE_URL"`
	RedisAddr string `mapstructure:"REDIS_ADDR"`
	Port      string `mapstructure:"PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"` // Wajib ada untuk Auth User
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../")
	viper.AddConfigPath("/app")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Info: .env file not found, relying on System Environment Variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// Manual fallback ke Environment Variables jika Viper gagal/kosong
	if config.DBUrl == "" {
		config.DBUrl = os.Getenv("DATABASE_URL")
	}
	if config.RedisAddr == "" {
		config.RedisAddr = os.Getenv("REDIS_ADDR")
	}
	if config.Port == "" {
		config.Port = os.Getenv("PORT")
	}
	if config.JWTSecret == "" {
		config.JWTSecret = os.Getenv("JWT_SECRET")
	}

	// Validasi Critical Config
	if config.DBUrl == "" {
		log.Fatal("FATAL: DATABASE_URL is empty. Please check your docker-compose.yml or .env")
	}
	if config.JWTSecret == "" {
		log.Fatal("FATAL: JWT_SECRET is empty. Auth cannot work without it.")
	}

	return &config
}