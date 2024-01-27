package dto

type UserRegisterIn struct {
	Name     string
	Email    string
	Password string
}

type UserActivateIn struct {
	Code string
}

type UserAuthenticateIn struct {
	Email    string
	Password string
}

type UserAuthenticateOut struct {
	AccessToken string
}

type AccessTokenRefreshIn struct {
	AccessToken string
}

type AccessTokenRefreshOut struct {
	AccessToken string
}

type AccessTokenParseIn struct {
	AccessToken string
}

type ParseTokenOut struct {
	UserID int64
}
