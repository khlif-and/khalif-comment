//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"khalif-comment/internal/config"
	"khalif-comment/internal/domain"
	"khalif-comment/internal/handler"
	"khalif-comment/internal/repository"
	"khalif-comment/internal/usecase"

)

func InitializeApp() (*App, error) {
	wire.Build(
		// 1. Config & Providers
		config.LoadConfig,
		ProvideDB,
		ProvideRedis,

		// 2. Repositories
		repository.NewCommentRepository,
		repository.NewCacheRepository,

		// 3. Bind Interfaces ke Implementasi (Repository)
		wire.Bind(new(domain.CommentRepository), new(*repository.CommentRepo)),
		wire.Bind(new(domain.RedisRepository), new(*repository.RedisRepo)),

		// 4. UseCases
		usecase.NewCommentUseCase,

		// 5. Bind Interfaces ke Implementasi (UseCase)
		wire.Bind(new(domain.CommentUseCase), new(*usecase.CommentUC)),

		// 6. Handlers
		handler.NewCommentHandler,

		// 7. App Entry Point
		NewApp,
	)
	return &App{}, nil
}