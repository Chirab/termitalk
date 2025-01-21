package app

type Memory struct {
	token string

	serverSideToken string
}

func (am *Memory) GetAuthId() string {
	return am.token
}

func (am *Memory) SetAuthId(token string) {
	am.token = token
}

func (am *Memory) DecodeUserInfo(i string) {
	am.serverSideToken = i
}

func (am *Memory) GetServerSideToken() string {
	return am.serverSideToken
}
