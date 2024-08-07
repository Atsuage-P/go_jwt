package infrastructure

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_$GOFILE -package=mock

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go_jwt/domain"
	"go_jwt/domain/model"
	"go_jwt/infrastructure/sqlc"
	apperror "go_jwt/internal/errors"
)

// mockでテストするためにsqlcのQuerier interfaceをWrapする
type Querier interface {
	ExistsUser(ctx context.Context, email string) (bool, error)
	GetLastInsertID(ctx context.Context) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	InsertUser(ctx context.Context, arg sqlc.InsertUserParams) (sql.Result, error)
}

type userRepository struct {
	db Querier
}

func NewUserRepository(db sqlc.Querier) domain.UserRepository {
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNoUser
		}

		return nil, fmt.Errorf("DB Error: %w", err)
	}
	user := model.User{
		ID:       int64(result.UserID),
		Password: result.Password,
	}

	return &user, nil
}
