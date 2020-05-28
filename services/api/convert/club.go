package convert

import (
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/services/api/models"
)

type converterClub struct {
}

func NewConvertClub() *converterClub {
	return &converterClub{}
}

func (cc *converterClub) ConvertModelToStructDB(clubModel *models.Club) *structure.Club {
	clubStruct := &structure.Club{}
	clubStruct.Name = clubModel.Name
	clubStruct.Tagline = clubModel.Tagline
	clubStruct.Description = clubModel.Description
	clubStruct.Address = clubModel.Address
	clubStruct.Email = clubModel.Email
	clubStruct.Phone = clubModel.Phone

	return clubStruct
}
