package service

import (
	"go_jwt/internal/cert"
	"strings"
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	as := NewAuthService()
	hashedPassword, _ := as.HashPassword("password123")

	data := []struct {
		name           string
		password       string
		hashedPassword string
		isError        bool
	}{
		{
			name:           "正常系",
			password:       "password123",
			hashedPassword: hashedPassword,
			isError:        false,
		},
		{
			name:           "異常系_パスワードが空",
			password:       "",
			hashedPassword: hashedPassword,
			isError:        true,
		},
		{
			name:           "異常系_パスワード不一致",
			password:       "hoge",
			hashedPassword: hashedPassword,
			isError:        true,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := as.VerifyPassword(hashedPassword, d.password)
			if d.isError {
				// 異常系だがerrがnilのケース
				if err == nil {
					t.Errorf("case: %s", d.name)
				}
			} else {
				// 正常系だがerrが発生しているケース
				if err != nil {
					t.Errorf("case: %s, wantErr: <nil>, gotErr: %v", d.name, err)
				}
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	as := NewAuthService()

	data := []struct {
		name     string
		password string
		isError  bool
	}{
		{
			name:     "正常系",
			password: "password123",
			isError:  false,
		},
		{
			name:     "正常系_パスワード長さ上限以内",
			password: strings.Repeat("a", 72),
			isError:  false,
		},
		{
			name:     "異常系_パスワードが空",
			password: "",
			isError:  true,
		},
		{
			name:     "異常系_パスワード長さ上限以上",
			password: strings.Repeat("a", 73),
			isError:  true,
		},
	}

	for _, d := range data {
		hashedPassword, err := as.HashPassword(d.password)
		if d.isError {
			// 異常系だがerrがnilのケース
			if err == nil {
				t.Errorf("case: %s", d.name)
			}
		} else {
			// 正常系だがerrが発生しているケース
			if err != nil {
				t.Errorf("case: %s, wantErr: <nil>, got: %s, gotErr: %v", d.name, hashedPassword, err)
			}
		}
	}
}

func TestCreateToken(t *testing.T) {
	as := NewAuthService()

	data := []struct {
		name    string
		userID  int64
		privKey []byte
		isError bool
	}{
		{
			name:    "正常系",
			userID:  1,
			privKey: cert.RawPrivKey,
			isError: false,
		},
	}

	for _, d := range data {
		_, err := as.CreateToken(d.userID)
		if d.isError {
			// 異常系だがerrがnilのケース
			if err == nil {
				t.Errorf("case: %s", d.name)
			}
		} else {
			// 正常系だがerrが発生しているケース
			if err != nil {
				t.Errorf("case: %s, wantErr: <nil>, gotErr: %v", d.name, err)
			}
		}
	}
}
