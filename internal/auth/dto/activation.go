package dto

type UserActivationEmailMessage struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
