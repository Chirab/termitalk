package app

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(usr string, passwd string) bool {
	return true
}

func (a *Auth) IsAuth() bool {
	return true
}
