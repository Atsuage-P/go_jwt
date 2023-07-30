package infrastructure

import (
	"context"
	"go_oauth/domain"
	"go_oauth/domain/model"
)

type userRepository struct {
}

func NewUserRepository() domain.UserRepository {
	return &userRepository{}
}

func (r *userRepository) CreateUser(ctx context.Context, signupInfo model.SignupInput) error {
	return nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}