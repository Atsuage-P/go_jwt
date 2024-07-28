package registry

import (
	"go_jwt/application/usecase"
	"go_jwt/domain/service"
	"go_jwt/env"
	"go_jwt/infrastructure"
	"go_jwt/infrastructure/config"
	"go_jwt/infrastructure/sqlc"
	"go_jwt/presentation/handler"
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
