package club

import (
	"fmt"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/transport/files"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type DeleteHandler struct {
	handler.HandlerBase
	id       int64
	response *handler.Response
}

func NewDeleteHandler() handler.IHandler {
	return &DeleteHandler{}
}

func (d *DeleteHandler) BindContext(ctx *gin.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	d.id = id
	return nil
}

func (d *DeleteHandler) Validate() *validators.Validate {
	return nil
}

func (d *DeleteHandler) Process() *handler.Response {
	clubRepo := db.NewClubRepository()
	clubDB, err := clubRepo.FindByID(d.id)
	if err != nil {
		return d.response.Error(500, errors.Wrap(err, "fail find club"))
	}

	if clubDB == nil {
		return d.response.Error(400, errors.New("club not found"))
	}

	// remove clubs images
	var removeImageErrors []string
	cdn := files.GetCDNStorage()
	if clubDB.Logo != "" {
		err := cdn.DELETE(clubDB.Logo)
		if err != nil {
			removeImageErrors = append(removeImageErrors, "fail remove club logo: " + err.Error())
		} else {
			clubDB.Logo = ""
		}
	}

	if clubDB.CoverImage != "" {
		err := cdn.DELETE(clubDB.CoverImage)
		if err != nil {
			removeImageErrors = append(removeImageErrors, "fail remove club cover: " + err.Error())
		} else {
			clubDB.CoverImage = ""
		}
	}

	if len(removeImageErrors) > 0 {
		logger.Log.Error(errors.New(fmt.Sprintf("fail remove club images. Club ID: %d, errors: %s", clubDB.ID, strings.Join(removeImageErrors, ","))))
	}

	err = clubRepo.Delete(clubDB.ID)
	if err != nil {
		return d.response.Error(500, errors.Wrap(err, "fail delete club"))
	}

	return d.response.Success("deleted")
}
