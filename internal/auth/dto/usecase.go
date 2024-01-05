package dto

type RegisterIn struct {
	Name     string
	Email    string
	Password string
}

type ActivateIn struct {
	Code string
}

type AuthenticateIn struct {
	Email    string
	Password string
}

type AuthenticateOut struct {
	AccessToken string
}

type RefreshTokenIn struct {
	AccessToken string
}

type RefreshTokenOut struct {
	AccessToken string
}

type ParseTokenIn struct {
	AccessToken string
}

type ParseTokenOut struct {
	UserID int64
}
