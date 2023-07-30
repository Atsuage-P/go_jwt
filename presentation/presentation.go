package presentation

type AuthHandler interface {
	SignUp()
	Login()
	Logout()
}