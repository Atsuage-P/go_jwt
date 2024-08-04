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
						model.SignupInput{
							UserName: "user1",
							Email:    "test1@mail.com",
							Password: "password",
						},
					).Return(int64(1), nil)

					return m
				},
				authService: func(t *testing.T) domain.AuthService {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockAuthService(ctrl)
					m.EXPECT().HashPassword(gomock.Any()).Return("password", nil)
					m.EXPECT().CreateToken(int64(1)).Return("token1", nil)

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
