//nolint:funlen
package usecase_test

import (
	"context"
	"errors"
	"go_jwt/application/usecase"
	"go_jwt/domain"
	"go_jwt/domain/model"
	"go_jwt/mock"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func Test_auth_SignUp(t *testing.T) {
	type fields struct {
		userRepository func(*testing.T) domain.UserRepository
		authService    func(*testing.T) domain.AuthService
	}
	type args struct {
		ctx      context.Context
		username string
		email    string
		password string
	}

	signupInput := model.SignupInput{
		UserName: "user1",
		Email:    "test1@mail.com",
		Password: "password",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.SignupOutput
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(false, nil)

					m.EXPECT().CreateUser(
						gomock.Any(),
						signupInput,
					).Return(int64(1), nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().HashPassword(gomock.Any()).Return("password", nil)
					m.EXPECT().CreateToken(gomock.Any()).Return("token1", nil)

					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want: &model.SignupOutput{
				Token: "token1",
			},
		},
		{
			name: "異常系_ユーザー検索でエラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(false, errors.New("error"))

					return m
				},
				authService: func(*testing.T) domain.AuthService {
					return nil
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_ユーザー検索でユーザー重複エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(true, nil)

					return m
				},
				authService: func(*testing.T) domain.AuthService {
					return nil
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_パスワードハッシュ化エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(false, nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().HashPassword(gomock.Any()).Return("", errors.New("error"))

					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_ユーザー作成エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(false, nil)

					m.EXPECT().CreateUser(
						gomock.Any(),
						signupInput,
					).Return(int64(1), errors.New("error"))

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().HashPassword(gomock.Any()).Return("password", nil)

					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_トークン生成エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().ExistsUser(
						gomock.Any(),
						gomock.Any(),
					).Return(false, nil)

					m.EXPECT().CreateUser(
						gomock.Any(),
						signupInput,
					).Return(int64(1), nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().HashPassword(gomock.Any()).Return("password", nil)
					m.EXPECT().CreateToken(gomock.Any()).Return("", errors.New("error"))

					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "user1",
				email:    "test1@mail.com",
				password: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := usecase.NewAuthUsecase(tt.fields.userRepository(t), tt.fields.authService(t))
			got, err := a.SignUp(tt.args.ctx, tt.args.username, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.SignUp() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.SignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Login(t *testing.T) {
	type fields struct {
		userRepository func(*testing.T) domain.UserRepository
		authService    func(*testing.T) domain.AuthService
	}
	type args struct {
		ctx           context.Context
		email         string
		inputPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.LoginOutput
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().FindUserByEmail(
						gomock.Any(),
						gomock.Any(),
					).Return(&model.User{}, nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().VerifyPassword(
						gomock.Any(),
						gomock.Any(),
					).Return(nil)
					m.EXPECT().CreateToken(gomock.Any()).Return("token", nil)

					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				email:         "test1@mail.com",
				inputPassword: "password",
			},
			want: &model.LoginOutput{
				Token: "token",
			},
		},
		{
			name: "異常系_メールによるユーザー検索でエラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().FindUserByEmail(
						gomock.Any(),
						gomock.Any(),
					).Return(nil, errors.New("error"))

					return m
				},
				authService: func(*testing.T) domain.AuthService {
					return nil
				},
			},
			args: args{
				ctx:           context.Background(),
				email:         "test1@mail.com",
				inputPassword: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_パスワード検証で内部エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().FindUserByEmail(
						gomock.Any(),
						gomock.Any(),
					).Return(&model.User{}, nil)

					return m
				},
				authService: func(*testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().VerifyPassword(
						gomock.Any(),
						gomock.Any(),
					).Return(errors.New("error"))

					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				email:         "test1@mail.com",
				inputPassword: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_パスワード検証で不一致エラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().FindUserByEmail(
						gomock.Any(),
						gomock.Any(),
					).Return(&model.User{}, nil)

					return m
				},
				authService: func(*testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().VerifyPassword(
						gomock.Any(),
						gomock.Any(),
					).Return(bcrypt.ErrMismatchedHashAndPassword)

					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				email:         "test1@mail.com",
				inputPassword: "password",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_トークン生成でエラー",
			fields: fields{
				userRepository: func(t *testing.T) domain.UserRepository {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockUserRepository(ctrl)
					m.EXPECT().FindUserByEmail(
						gomock.Any(),
						gomock.Any(),
					).Return(&model.User{}, nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().VerifyPassword(
						gomock.Any(),
						gomock.Any(),
					).Return(nil)
					m.EXPECT().CreateToken(gomock.Any()).Return("", errors.New("error"))

					return m
				},
			},
			args: args{
				ctx:           context.Background(),
				email:         "test1@mail.com",
				inputPassword: "password",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := usecase.NewAuthUsecase(tt.fields.userRepository(t), tt.fields.authService(t))
			got, err := a.Login(tt.args.ctx, tt.args.email, tt.args.inputPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.Login() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Hello(t *testing.T) {
	type fields struct {
		userRepository func(*testing.T) domain.UserRepository
		authService    func(*testing.T) domain.AuthService
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.APIOutput
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				userRepository: func(*testing.T) domain.UserRepository {
					return nil
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().InvalidateToken(gomock.Any()).Return(nil)

					return m
				},
			},
			args: args{token: "test"},
			want: &model.APIOutput{Message: "success"},
		},
		{
			name: "異常系_Tokenが無効",
			fields: fields{
				userRepository: func(*testing.T) domain.UserRepository {
					return nil
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().InvalidateToken(gomock.Any()).Return(errors.New("error"))

					return m
				},
			},
			args:    args{token: "test"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_Tokenが空",
			fields: fields{
				userRepository: func(*testing.T) domain.UserRepository {
					return nil
				},
				authService: func(*testing.T) domain.AuthService {
					return nil
				},
			},
			args:    args{token: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := usecase.NewAuthUsecase(tt.fields.userRepository(t), tt.fields.authService(t))
			got, err := a.Hello(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.Hello() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Hello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Logout(t *testing.T) {
	type fields struct {
		userRepository func(t *testing.T) domain.UserRepository
		authService    func(t *testing.T) domain.AuthService
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				userRepository: func(*testing.T) domain.UserRepository {
					return nil
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().InvalidateToken(gomock.Any()).Return(nil)

					return m
				},
			},
			args: args{
				token: "token",
			},
		},
		{
			name: "異常系_無効なトークン",
			fields: fields{
				userRepository: func(*testing.T) domain.UserRepository {
					return nil
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().InvalidateToken(gomock.Any()).Return(errors.New("error"))

					return m
				},
			},
			args: args{
				token: "token",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := usecase.NewAuthUsecase(tt.fields.userRepository(t), tt.fields.authService(t))
			if err := a.Logout(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("auth.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
