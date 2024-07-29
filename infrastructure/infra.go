package infrastructure

import (
	"context"
	"fmt"
	"go_jwt/domain"
	"go_jwt/domain/model"
	"go_jwt/infrastructure/sqlc"
)

type userRepository struct {
	db *sqlc.Queries
}

func NewUserRepository(db *sqlc.Queries) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, signupInfo model.SignupInput) (int64, error) {
	args := sqlc.InsertUserParams{
		UserName: signupInfo.UserName,
		Email:    signupInfo.Email,
		Password: signupInfo.Password,
	}
	result, err := r.db.InsertUser(ctx, args)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *userRepository) ExistsUser(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.ExistsUser(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	result, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("DB Error: %w", err)
	}
	user := model.User{
		ID:       int64(result.UserID),
		Password: result.Password,
	}

	return &user, nil
}
