package models

import (
	"net/http"
)

type UserSignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a UserLoginRequest) Bind(r *http.Request) error {
	return nil
}

func (l UserSignupRequest) Bind(r *http.Request) error {
	return nil
}

type TokenResponse struct {
	Token string `json:"token"`
}
