package model

type SignupInput struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupOutput struct {
	Token string `json:"token"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type APIInput struct {
	Token string `json:"token"`
}

type APIOutput struct {
	Message string `json:"message"`
}
