package usecase

import (
	"context"
	"errors"
	"fmt"
	"go_jwt/application"
	"go_jwt/domain"
	"go_jwt/domain/model"

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
		return nil, fmt.Errorf("User Find Error: %v", err)
	}
	// ユーザーがいれば早期リターン
	if exists {
		// TODO: errorsに移す
		return nil, errors.New("このメールアドレスは既に登録されています。別のメールアドレスを使用してください。")
	}

	// パスワードを暗号化
	hashedPassword, err := a.authService.HashPassword(ctx, password)
	if err != nil {
		return nil, fmt.Errorf("Password Hashed Error: %v", err)
	}
	signupInfo := model.SignupInput{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	// ユーザーのDB登録
	userID, err := a.userRepository.CreateUser(ctx, signupInfo)
	if err != nil {
		return nil, fmt.Errorf("DB Error: %v", err)
	}
	token, err := a.authService.CreateToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Create Token Error: %v", err)
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
		// TODO: エラー処理
		return nil, err
	}
	if user == nil {
		// TODO: エラー処理
		return nil, errors.New("ユーザーが見つかりません。")
	}

	// パスワードの検証
	if err := a.authService.VerifyPassword(ctx, user.Password, inputPassword); err == bcrypt.ErrMismatchedHashAndPassword {
		// TODO: errorsに移す
		return nil, errors.New("メールアドレスまたはパスワードが間違っています。")
	} else if err != nil {
		// InternalError
		return nil, err
	}
	token, err := a.authService.CreateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	output := model.LoginOutput{
		Token: token,
	}
	return &output, nil
}

func (a *auth) Logout(ctx context.Context, token string) error {
	if err := a.authService.InvalidateToken(ctx, token); err != nil {
		return err
	}
	return nil
}

func (a *auth) Hello(ctx context.Context, token string) (*model.APIOutput, error) {
	if err := a.authService.InvalidateToken(ctx, token); err != nil {
		return nil, err
	}
	output := model.APIOutput{
		Message: "success",
	}

	return &output, nil
}
