package application

import (
	"context"
	"go_oauth/domain/model"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, username, email, password string) (*model.SignupOutput, error)
	Login(ctx context.Context, email, password string) (*model.LoginOutput, error)
	Logout(ctx context.Context, token string) error
}
