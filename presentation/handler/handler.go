package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go_jwt/application"
	apperror "go_jwt/internal/errors"
	"go_jwt/internal/logger"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase application.AuthUsecase
}

func NewAuthHandler(authUsecase application.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (ah *AuthHandler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger()
	ctx := logger.WithContext(r.Context(), log)

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic",
				zap.Any("panic_msg", r),
			)
			http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		}
	}()

	var data struct {
		UserName string `json:"user_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"message": "Bad Request"}`, http.StatusBadRequest)

		return
	}

	res, err := ah.authUsecase.SignUp(ctx, data.UserName, data.Email, data.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrDuplicateID) {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, `{"message": "Bad Request"}`, http.StatusBadRequest)
			log.Error(err.Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
			log.Error(err.Error())
		}

		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}

func (ah *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger()
	ctx := logger.WithContext(r.Context(), log)

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic",
				zap.Any("panic_msg", r),
			)
			http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		}
	}()

	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, `{"message": "Bad Request"}`, http.StatusBadRequest)
		log.Error(err.Error())

		return
	}

	res, err := ah.authUsecase.Login(ctx, data.Email, data.Password)
	if err != nil {
		if errors.Is(err, apperror.ErrWrongLoginInfo) {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		} else if errors.Is(err, apperror.ErrNoUser) {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, `"message": "Bad Request"`, http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		}
		log.Error(err.Error())

		return
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}

// func (ah *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
// }

func (ah *AuthHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.NewLogger()

	defer func() {
		if r := recover(); r != nil {
			log.Error("panic",
				zap.Any("panic_msg", r),
			)
			http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		}
	}()

	var data struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Error(err.Error())

		return
	}

	res, err := ah.authUsecase.Hello(data.Token)
	if err != nil {
		if errors.Is(err, apperror.ErrTokenIsNone) {
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, `{"message": "Unauthorized"}`, http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Error(err.Error())

		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
