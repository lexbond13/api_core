package models

import (
	"github.com/lexbond13/api_core/services/api/handler/validators"
)

type Club struct {
	Name        string `json:"name"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Logo        string `json:"logo"`
	CoverImage  string `json:"cover_image"`
}

// Validate
func (ad *Club) Validate() *validators.Validate {
	vErrs := &validators.Validate{}
	vErrs.Required("name", ad.Name)
	vErrs.Required("description", ad.Description)
	vErrs.Required("address", ad.Address)
	vErrs.Required("email", ad.Email)
	vErrs.Required("phone", ad.Phone)

	return vErrs
}
