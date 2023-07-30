package registry

import (
	"go_oauth/application/usecase"
	"go_oauth/domain/service"
	"go_oauth/infrastructure"
	"go_oauth/presentation/handler"
)

func AuthRegistry() handler.AuthHandler {
	userRepository := infrastructure.NewUserRepository()
	authService := service.NewAuthService()
	authUsecase := usecase.NewAuthUsecase(userRepository, authService)
	authHandler := handler.NewAuthHandler(authUsecase)
	return *authHandler
}