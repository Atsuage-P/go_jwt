package domain

import (
	"context"
	"go_oauth/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, signupInfo model.SignupInput) error
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	VerifyPassoword(ctx context.Context, password, hashedPassword string)
	HashPassowrd(ctx context.Context, password string) (string, error)
	CreateToken(ctx context.Context, userID int) (string, error)
	InvalidateToken(ctx context.Context, token string) error
}
