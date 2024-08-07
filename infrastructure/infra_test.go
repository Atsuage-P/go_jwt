//nolint:funlen
package infrastructure_test

import (
	"context"
	"database/sql"
	"errors"
	"go_jwt/domain/model"
	"go_jwt/infrastructure"
	"go_jwt/infrastructure/sqlc"
	"go_jwt/mock"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/mock/gomock"
)

func Test_userRepository_CreateUser(t *testing.T) {
	type fields struct {
		db func(*testing.T) infrastructure.Querier
	}
	type args struct {
		ctx        context.Context
		signupInfo model.SignupInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(sqlmock.NewResult(int64(1), int64(1)), nil)

					return m
				},
			},
			args: args{
				ctx: context.Background(),
				signupInfo: model.SignupInput{
					UserName: "user1",
					Email:    "test@mail.com",
					Password: "password",
				},
			},
			want: 1,
		},
		{
			name: "異常系_ユーザーインサートでエラー",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

					return m
				},
			},
			args:    args{},
			want:    0,
			wantErr: true,
		},
		{
			name: "異常系_インサート結果の取得でエラー",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(sqlmock.NewErrorResult(errors.New("error")), nil)

					return m
				},
			},
			args:    args{},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := infrastructure.NewUserRepository(tt.fields.db(t))
			got, err := r.CreateUser(tt.args.ctx, tt.args.signupInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("userRepository.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_ExistsUser(t *testing.T) {
	type fields struct {
		db func(*testing.T) infrastructure.Querier
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "正常系_ユーザーの重複なし",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().ExistsUser(gomock.Any(), gomock.Any()).Return(false, nil)

					return m
				},
			},
			args: args{
				ctx:   context.Background(),
				email: "test@mail.com",
			},
			want: false,
		},
		{
			name: "正常系_ユーザーの重複あり",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().ExistsUser(gomock.Any(), gomock.Any()).Return(true, nil)

					return m
				},
			},
			args: args{
				ctx:   context.Background(),
				email: "test@mail.com",
			},
			want: true,
		},
		{
			name: "異常系_ユーザー重複検索でエラー",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().ExistsUser(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))

					return m
				},
			},
			args:    args{},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := infrastructure.NewUserRepository(tt.fields.db(t))
			got, err := r.ExistsUser(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.ExistsUser() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("userRepository.ExistsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_FindUserByEmail(t *testing.T) {
	type fields struct {
		db func(*testing.T) infrastructure.Querier
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "正常系",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(sqlc.GetUserByEmailRow{
						UserID:   uint32(1),
						UserName: "test",
						Password: "password",
					}, nil)

					return m
				},
			},
			args: args{
				ctx:   context.Background(),
				email: "test@email.com",
			},
			want: &model.User{
				ID:       int64(1),
				Password: "password",
			},
		},
		{
			name: "異常系_ユーザーが見つからなかった",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(sqlc.GetUserByEmailRow{}, sql.ErrNoRows)

					return m
				},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系_ユーザー検索時にエラー",
			fields: fields{
				db: func(t *testing.T) infrastructure.Querier {
					t.Helper()
					ctrl := gomock.NewController(t)
					m := mock.NewMockQuerier(ctrl)
					m.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Return(sqlc.GetUserByEmailRow{}, errors.New("error"))

					return m
				},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := infrastructure.NewUserRepository(tt.fields.db(t))
			got, err := r.FindUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.FindUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
