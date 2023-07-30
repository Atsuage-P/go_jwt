package service

import (
	"context"
	"go_oauth/domain"
)

type authService struct {
}

func NewAuthService() domain.AuthService {
	return &authService{}
}

func (as *authService) VerifyPassoword(ctx context.Context, password, hashedPassword string) {

}

func (as *authService) HashPassowrd(ctx context.Context, password string) (string, error) {
 return "", nil
}

func (as *authService) CreateToken(ctx context.Context, userID int) (string, error) {
	return "", nil
}

func (as *authService) InvalidateToken(ctx context.Context, token string) error {
	return nil
}
