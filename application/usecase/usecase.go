package usecase

import (
	"context"
	"errors"
	"fmt"
	"go_jwt/application"
	"go_jwt/domain"
	"go_jwt/domain/model"
	apperror "go_jwt/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

type auth struct {
	userRepository domain.UserRepository
	authService    domain.AuthService
}

// AuthUsecaseインターフェースを満たす構造体ポインタを返す
func NewAuthUsecase(
	userRepository domain.UserRepository,
	authService domain.AuthService,
) application.AuthUsecase {
	return &auth{
		userRepository: userRepository,
		authService:    authService,
	}
}

func (a *auth) SignUp(ctx context.Context, username, email, password string) (*model.SignupOutput, error) {
	// 同Emailのユーザーの検索
	exists, err := a.userRepository.ExistsUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("User Find Error: %w", err)
	}
	// ユーザーがいれば早期リターン
	if exists {
		return nil, apperror.ErrDuplicateID
	}

	// パスワードを暗号化
	hashedPassword, err := a.authService.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("Password Hashed Error: %w", err)
	}
	signupInfo := model.SignupInput{
		UserName: username,
		Email:    email,
		Password: hashedPassword,
	}

	// ユーザーのDB登録
	userID, err := a.userRepository.CreateUser(ctx, signupInfo)
	if err != nil {
		return nil, fmt.Errorf("DB Error: %w", err)
	}
	token, err := a.authService.CreateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("Create Token Error: %w", err)
	}
	output := model.SignupOutput{
		Token: token,
	}

	return &output, nil
}

func (a *auth) Login(ctx context.Context, email, inputPassword string) (*model.LoginOutput, error) {
	// DB検索
	user, err := a.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// パスワードの検証
	err = a.authService.VerifyPassword(user.Password, inputPassword)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		// bcryptError
		return nil, apperror.ErrWrongLoginInfo
	} else if err != nil {
		// InternalError
		return nil, err
	}
	token, err := a.authService.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}
	output := model.LoginOutput{
		Token: token,
	}

	return &output, nil
}

func (a *auth) Logout(token string) error {
	if err := a.authService.InvalidateToken(token); err != nil {
		return err
	}

	return nil
}

func (a *auth) Hello(token string) (*model.APIOutput, error) {
	if err := a.authService.InvalidateToken(token); err != nil {
		return nil, err
	}
	output := model.APIOutput{
		Message: "success",
	}

	return &output, nil
}
