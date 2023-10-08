package handler

import (
	"encoding/json"
	"net/http"

	"go_oauth/application"
)

type AuthHandler struct {
	authUsecase application.AuthUsecase
}

func NewAuthHandler(authUsecase application.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// TODO: 入力値のバリデーション
func (ah *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	ctx := r.Context()
	res, err := ah.authUsecase.SignUp(ctx, data.UserName, data.Email, data.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
}

func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
}
