package app

type AuthMemory struct {
	token string

	serverSideToken string
}

func (am *AuthMemory) GetAuthId() string {
	return am.token
}

func (am *AuthMemory) SetAuthId(token string) {
	am.token = token
}

func (am *AuthMemory) DecodeUserInfo(i string) {
	am.serverSideToken = i
}

func (am *AuthMemory) GetServerSideToken() string {
	return am.serverSideToken
}
