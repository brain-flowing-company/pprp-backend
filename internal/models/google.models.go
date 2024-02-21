package models

type GoogleExchangeTokens struct {
	Authuser string
	Code     string
	Prompt   string
	Scope    string
	State    string
}

type GoogleUserInfo struct {
	Email string
}
