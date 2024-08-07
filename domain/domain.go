package domain

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOFILE -package=mock

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
	VerifyPassword(password, hashedPassword string) error
	HashPassword(password string) (string, error)
	CreateToken(userID int64) (string, error)
	InvalidateToken(token string) error
}
