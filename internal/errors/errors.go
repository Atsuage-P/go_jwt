package apperror

import (
	"errors"
	"fmt"
)

type AuthError struct {
	Message string
}

var (
	ErrNoUser         = errors.New("ユーザーが見つかりません")
	ErrDuplicateID    = errors.New("このメールアドレスは既に登録されています。別のメールアドレスを使用してください")
	ErrWrongLoginInfo = errors.New("メールアドレスまたはパスワードが間違っています")
	ErrPasswordIsNone = errors.New("パスワードが空です")
	ErrTokenIsNone    = errors.New("トークンが空です")
)

func (a *AuthError) Error() string {
	return fmt.Sprintf("Error: %s", a.Message)
}
