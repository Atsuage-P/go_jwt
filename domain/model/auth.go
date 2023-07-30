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
