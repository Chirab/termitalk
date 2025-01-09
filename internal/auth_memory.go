package app

type AuthMemory struct {
	authId string
}

func (am *AuthMemory) GetAuthId() string {
	return am.authId
}

func (am *AuthMemory) SetAuthId(authId string) {
	am.authId = authId
}
