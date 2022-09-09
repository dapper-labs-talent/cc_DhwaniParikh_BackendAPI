package models

import "net/http"

type UserUpdateRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (a UserUpdateRequest) Bind(r *http.Request) error {
	return nil

}
