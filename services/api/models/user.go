package models

import (
	"github.com/lexbond13/api_core/services/api/handler/validators"
	jwtGo "github.com/dgrijalva/jwt-go"
)

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type UserSessionData struct {
	ID    int64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Status string `json:"status"`
}

type Claims struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	jwtGo.StandardClaims
}

// Validate
func (u *User) Validate() *validators.Validate {
	vErrs := &validators.Validate{}
	vErrs.Required("Email", u.Email)
	vErrs.Required("Name", u.Name)
	vErrs.IsEmail(u.Email)

	return vErrs
}
