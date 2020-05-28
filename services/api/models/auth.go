package models

import "github.com/lexbond13/api_core/services/api/handler/validators"

type Auth struct {
	Email string `json:"email"`
}

type AuthActivate struct {
	Email string `json:"email"`
	ActiveKey string `json:"active_key"`
}

// Validate
func (a *Auth) Validate() *validators.Validate {
	vErrs := &validators.Validate{}
	vErrs.Required("email", a.Email)
	vErrs.IsEmail(a.Email)

	return vErrs
}

// Validate
func (a *AuthActivate) Validate() *validators.Validate {
	vErrs := &validators.Validate{}
	vErrs.Required("email", a.Email)
	vErrs.Required("active_key", a.ActiveKey)
	vErrs.IsEmail(a.Email)

	return vErrs
}
