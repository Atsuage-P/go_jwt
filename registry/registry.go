package registry

import (
	"go_oauth/application/usecase"
	"go_oauth/domain/service"
	"go_oauth/env"
	"go_oauth/infrastructure"
	"go_oauth/infrastructure/config"
	"go_oauth/infrastructure/sqlc"
	"go_oauth/presentation/handler"
)

func AuthRegistry(cnf *env.EnvConfig) handler.AuthHandler {
	db := config.ConnectDB(&cnf.DB)
	queries := sqlc.New(db)

	userRepository := infrastructure.NewUserRepository(queries)
	authService := service.NewAuthService()
	authUsecase := usecase.NewAuthUsecase(userRepository, authService)
	authHandler := handler.NewAuthHandler(authUsecase)
	return *authHandler
}
