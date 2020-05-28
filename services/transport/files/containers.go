package files

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/lexbond13/api_core/util"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"
)

type IFileContainer interface {
	Reader() (io.Reader, error)
	Extension() string
	FileName() string
	FileExt() string
	FileSize() int64
}

type BinDataFileContainer struct {
	bufDecodedData  *bytes.Buffer
	allowExtensions []string
	allowFileSize   int64
	fileName       string
	fileExt        string
	fileSize       int64
}

type Base64ImageContainer struct {
	readerBase64   io.Reader
	bufDecodedData *bytes.Buffer
	imageDecode image.Image

	allowExtensions []string
	allowFileSize   int64
	fileName       string
	fileExt        string
	fileSize       int64
}

// NewBinDataFileContainer
func NewBinDataFileContainer(r io.Reader, allowExtensions []string, allowFileSize int64) (*BinDataFileContainer, error) {
	var err error

	if r == nil {
		err = errors.New("reader not set")
		return nil, err
	}

	var Buf bytes.Buffer
	_, err = io.Copy(&Buf, r)

	if err != nil {
		err = errors.New("can't copy file bytes to Buffer")
		return nil, err
	}

	container := &BinDataFileContainer{bufDecodedData: &Buf}
	container.allowExtensions = allowExtensions
	container.allowFileSize = allowFileSize

	return container, nil
}

// Reader
func (bi *BinDataFileContainer) Reader() (io.Reader, error) {
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, bi.bufDecodedData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to copy data to reader")
	}

	tee := io.TeeReader(buf, bi.bufDecodedData)

	return tee, nil
}

// Extension
func (bi *BinDataFileContainer) Extension() string {
	if bi.fileExt != "" {
		return bi.fileExt
	}

	return ""
}

// SetFileName
func (bi *BinDataFileContainer) SetFileName(fileName string) {
	bi.fileName = fileName
}

// SetFileExt
func (bi *BinDataFileContainer) SetFileExt(fileExt string) {
	bi.fileExt = fileExt
}

// SetFileSize
func (bi *BinDataFileContainer) SetFileSize(fileSize int64) {
	bi.fileSize = fileSize
}

func (bi *BinDataFileContainer) FileName() string {
	return bi.fileName
}

func (bi *BinDataFileContainer) FileExt() string {
	return bi.fileExt
}

func (bi *BinDataFileContainer) FileSize() int64 {
	return bi.fileSize
}

// Validate
func (bi *BinDataFileContainer) Validate() error {
	var errs []string
	if err := validateSize(bi.fileSize, bi.allowFileSize); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateExtension(bi.fileExt, bi.allowExtensions); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}
	return nil
}

// NewBase64ImageContainer
func NewBase64ImageContainer(r io.Reader) (*Base64ImageContainer, error) {
	var err error

	if r == nil {
		err = errors.New("error: not set reader")
		return nil, err
	}
	b := &bytes.Buffer{}
	_, err = io.Copy(b, r)
	if err != nil {
		err = errors.Wrap(err, "Failed to decoded base64 image")

		return nil, err
	}
	bi := &Base64ImageContainer{bufDecodedData: &bytes.Buffer{}}
	base64String := strings.TrimRight(b.String(), "=")
	bytesImage, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(base64String)
	if err != nil {
		err = errors.Wrap(err, "Failed to decoded base64 image")

		return nil, err
	}
	decodedData := bytes.NewReader(bytesImage)
	tee := io.TeeReader(decodedData, bi.bufDecodedData)
	bi.imageDecode, bi.fileExt, err = image.Decode(tee)
	if err != nil {
		err = errors.Wrap(err, "Failed to get decoded base64 image")

		return nil, err
	}

	return bi, nil
}

// Reader
func (bc *Base64ImageContainer) Reader() (io.Reader, error) {

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, bc.bufDecodedData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to copy data to reader")
	}

	tee := io.TeeReader(buf, bc.bufDecodedData)

	return tee, nil
}

// Extension
func (bc *Base64ImageContainer) Extension() string {
	var err error

	if bc.fileExt != "" {

		return bc.fileExt
	}
	if bc.readerBase64 != nil {
		_, bc.fileExt, err = image.DecodeConfig(bc.readerBase64)
		if err != nil {

			return ""
		}

		return bc.fileExt
	}

	return ""
}

//Image this method need for test
func (bc *Base64ImageContainer) Image() (img image.Image, name string, ext string, err error) {
	return bc.imageDecode, bc.fileName, bc.fileExt, nil
}

func (bc *Base64ImageContainer) FileName() string {
	return bc.fileName
}

func (bc *Base64ImageContainer) FileExt() string {
	return bc.fileExt
}

func (bc *Base64ImageContainer) FileSize() int64 {
	return bc.fileSize
}

// Validate
func Validate(container IFileContainer, allowExtensions []string, allowFileSize int64) error {
	var errs []string
	if container.FileSize() > allowFileSize {
		errs = append(errs, fmt.Sprintf("max image size allow: %dMB.", util.ConvertCountBytesToCountMegabytes(allowFileSize)))
	}

	if len(allowExtensions) > 0 {
		for _, allowExt := range allowExtensions {
			if container.FileExt() == strings.TrimSpace(allowExt) {
				return nil
			}
		}
		errs = append(errs, fmt.Sprintf("only extensions %s allowed.", strings.Join(allowExtensions, ",")))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}
	return nil
}

// TODO deprecated
func validateSize(size, allow int64) error {
	if size > allow {
		return errors.New(fmt.Sprintf("max image size allow: %dMB.", util.ConvertCountBytesToCountMegabytes(allow)))
	}
	return nil
}

func validateExtension(extension string, allow []string) error {
	if len(allow) == 0 {
		return nil
	}

	for _, allowExt := range allow {
		if extension == strings.TrimSpace(allowExt) {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("only extensions %s allowed.", strings.Join(allow, ",")))
}
