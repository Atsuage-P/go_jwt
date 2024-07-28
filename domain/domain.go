package domain

import (
	"context"
	"go_jwt/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, signupInfo model.SignupInput) (int64, error)
	ExistsUser(ctx context.Context, email string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	VerifyPassword(ctx context.Context, password, hashedPassword string) error
	HashPassword(ctx context.Context, password string) (string, error)
	CreateToken(ctx context.Context, userID int64) (string, error)
	InvalidateToken(ctx context.Context, token string) error
}
