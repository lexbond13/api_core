package club

import (
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/services/api/convert"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/api/models"
	"github.com/lexbond13/api_core/services/transport/files"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"strings"
)

type AddHandler struct {
	handler.HandlerBase
	request *models.Club
	response *handler.Response
	allowFileParams *config.UploadFiles
}

// NewAddHandler
func NewAddHandler(params *config.Params) handler.IHandler {
	return &AddHandler{
		allowFileParams: params.UploadFiles,
	}
}

// BindContext
func (a *AddHandler) BindContext(ctx *gin.Context) error {
	a.Ctx = ctx
	err := ctx.ShouldBindBodyWith(&a.request, binding.JSON)
	if err != nil {
		return err
	}

	return nil
}

// Validate
func (a *AddHandler) Validate() *validators.Validate {
	return a.request.Validate()
}

// Process
func (a *AddHandler) Process() *handler.Response {

	userIdentity, err := a.GetIdentity()
	if err != nil {
		return a.response.Error(401, errors.Wrap(err, "fail find user identity"))
	}

	// convert request data to db struct
	clubDB := convert.NewConvertClub().ConvertModelToStructDB(a.request)
	// uploading images
	cdnSender := files.GetCDNStorage()
	links, err := a.UploadImages(cdnSender)
	if err != nil {
		removeErr := a.RemoveImages(cdnSender, links["logo"], links["cover"])
		if removeErr != nil {
			logger.Log.Error(errors.Wrap(removeErr, "fail delete images"))
		}
		return a.response.Error(500, errors.Wrap(err, "fail upload images"))
	}

	clubDB.Logo = links["logo"]
	clubDB.CoverImage = links["cover"]
	clubDB.OwnerID = userIdentity.ID
	clubDB.Status = structure.ClubStatusNew

	clubRepo := db.NewClubRepository()
	err = clubRepo.Create(clubDB)
	if err != nil {
		removeErr := a.RemoveImages(cdnSender, clubDB.Logo, clubDB.CoverImage)
		if removeErr != nil {
			logger.Log.Error(errors.Wrap(removeErr, "fail delete images"))
		}
		return a.response.Error(500, errors.Wrap(err, "fail create club"))
	}

	return a.response.Success(clubDB)
}

// UploadImages
func (a *AddHandler) UploadImages(cdn files.ICDNFileStorage) (map[string]string, error) {

	links := make(map[string]string)
	var uploadFileErrors []string

	if a.request.Logo != "" {
		logoContainer, err := files.NewBase64ImageContainer(strings.NewReader(a.request.Logo))
		if err != nil {
			uploadFileErrors = append(uploadFileErrors, err.Error())
		} else {
			err := files.Validate(logoContainer, a.allowFileParams.AllowExtensions, a.allowFileParams.AllowSize)
			if err != nil {
				uploadFileErrors = append(uploadFileErrors, err.Error())
			} else {
				logoURL, err := cdn.PUT(logoContainer)
				if err != nil {
					uploadFileErrors = append(uploadFileErrors, err.Error())
				} else {
					links["logo"] = logoURL
				}
			}
		}
	}

	if a.request.CoverImage != "" {
		coverContainer, err := files.NewBase64ImageContainer(strings.NewReader(a.request.CoverImage))
		if err != nil {
			uploadFileErrors = append(uploadFileErrors, err.Error())
		} else {
			err := files.Validate(coverContainer, a.allowFileParams.AllowExtensions, a.allowFileParams.AllowSize)
			if err != nil {
				uploadFileErrors = append(uploadFileErrors, err.Error())
			} else {
				coverURL, err := cdn.PUT(coverContainer)
				if err != nil {
					uploadFileErrors = append(uploadFileErrors, err.Error())
				} else {
					links["cover"] = coverURL
				}
			}
		}
	}

	if len(uploadFileErrors) > 0 {
		return links, errors.New("fail upload images: " + strings.Join(uploadFileErrors, ","))
	}

	return links, nil
}

// RemoveImages
func (a *AddHandler) RemoveImages(cdn files.ICDNFileStorage, links ...string) error {
	var removeFileErrors []string
	for _, link := range links {
		if link == "" {
			continue
		}

		if err := cdn.DELETE(link); err != nil {
			removeFileErrors = append(removeFileErrors, err.Error())
		}
	}

	if len(removeFileErrors) > 0 {
		return errors.New("fail remove images: " + strings.Join(removeFileErrors, ","))
	}

	return nil
}
