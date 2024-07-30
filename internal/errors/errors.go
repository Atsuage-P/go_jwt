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
)

func (a *AuthError) Error() string {
	return fmt.Sprintf("Error: %s", a.Message)
}
