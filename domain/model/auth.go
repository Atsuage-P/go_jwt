package model

type SignupInput struct {
	Username string
	Email    string
	Password string
}

type SignupOutput struct {
	Token string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token string
}

type APIInput struct {
	Token string
}

type APIOutput struct {
	Message string
}
