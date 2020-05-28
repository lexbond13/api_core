package file

import (
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/db"
	"github.com/lexbond13/api_core/module/db/structure"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/lexbond13/api_core/services/api/handler/validators"
	"github.com/lexbond13/api_core/services/transport/files"
	"github.com/lexbond13/api_core/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
)

type UploadHandler struct {
	handler.HandlerBase
	File   *multipart.FileHeader
	params *config.UploadFiles
	response *handler.Response
}

func NewUploadHandler(params *config.Params) handler.IHandler {
	return &UploadHandler{
		params: params.UploadFiles,
	}
}

func (up *UploadHandler) BindContext(ctx *gin.Context) error {
	up.Ctx = ctx
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	up.File = file
	return nil
}

func (up *UploadHandler) Validate() *validators.Validate {
	return nil
}

// Process
func (up *UploadHandler) Process() *handler.Response {

	fileData, err := up.File.Open()
	if err != nil {
		return up.response.Error(500, errors.Wrap(err, "fail open uploaded file"))
	}

	defer func() {
		if err := fileData.Close(); err != nil {
			logger.Log.Error(errors.Wrap(err, "fail close file"))
		}
	}()

	fileContainer, err := files.NewBinDataFileContainer(fileData, up.params.AllowExtensions, up.params.AllowSize)
	if err != nil {
		return up.response.Error(500, err)
	}

	fileContainer.SetFileName(util.FileName(up.File.Filename))
	fileContainer.SetFileExt(util.FileExt(up.File.Filename))
	fileContainer.SetFileSize(up.File.Size)

	vErrs := fileContainer.Validate()
	if vErrs != nil {
		return up.response.Error(400, errors.New("fail validate: " + vErrs.Error()))
	}

	identity, err := up.GetIdentity()
	if err != nil {
		return up.response.Error(500, errors.Wrap(err, "fail to get user identity"))
	}

	cdnSender := files.GetCDNStorage()
	link, err := cdnSender.PUT(fileContainer)
	if err != nil {
		return up.response.Error(500, errors.Wrap(err, "fail send file to cdn storage"))
	}

	file := &structure.File{
		UserID:       identity.ID,
		OriginalName: up.File.Filename,
		ExternalURL:  link,
		Size:         fileContainer.FileSize(),
	}

	err = db.NewFileRepository().Create(file)
	if err != nil {
		return up.response.Error(500, errors.Wrap(err, "fail create file record to db"))
	}

	return up.response.Success(link)
}
