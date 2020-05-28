package files

import (
	"bytes"
	"fmt"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/module/logger"
	"github.com/lexbond13/api_core/util"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

var cdnStorage ICDNFileStorage

type ICDNFileStorage interface {
	PUT(file IFileContainer) (link string, err error)
	DELETE(link string) error
}

type SelCDN struct {
	URL            string
	AuthURL        string
	ContainerName  string
	AuthUser       string
	AuthKey        string
	accessToken    string
	timeOutRequest time.Duration
	devProdMode    bool
	isDebug        bool
}

func Init(config *config.SelCDN, isDebug bool) {
	cdnStorage = &SelCDN{
		URL:            config.URL,
		AuthURL:        config.AuthURL,
		ContainerName:  config.ContainerName,
		AuthUser:       config.User,
		AuthKey:        config.Key,
		timeOutRequest: 60,
		isDebug:        isDebug,
		devProdMode:    config.DevProdMode,
	}
}

// PUT send file to storage
func (sd *SelCDN) PUT(file IFileContainer) (link string, err error) {

	// choose container for uploading
	containerName := sd.ContainerName
	if sd.devProdMode {
		if sd.isDebug {
			containerName = "dev." + sd.ContainerName
		} else {
			containerName = "prod." + sd.ContainerName
		}
	}

	newFileName := fmt.Sprintf("%s_%d.%s", util.FilterOnlyCharsNums(file.FileName()), time.Now().Unix(), file.FileExt())
	cdnURL := fmt.Sprintf("%s/%s/%s", sd.URL, containerName, newFileName)

	fileBytes, _ := file.Reader()
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, fileBytes)

	req, err := http.NewRequest("PUT", cdnURL, bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		return "", err
	}

	accessToken, err := sd.getAccessToken()
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Auth-Token", accessToken)
	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{Timeout: sd.timeOutRequest * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Log.Error(errors.Wrap(err, "fail close file body"))
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", errors.New("fail upload file: " + resp.Status)
	}

	return cdnURL, nil
}

// DELETE
func (sd *SelCDN) DELETE(link string) error {
	req, err := http.NewRequest("DELETE", link, nil)
	if err != nil {
		return err
	}

	accessToken, err := sd.getAccessToken()
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", accessToken)
	req.Header.Set("Content-Type", "text/plain")

	client := http.Client{Timeout: sd.timeOutRequest * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Log.Error(errors.Wrap(err, "fail close file body"))
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return errors.New("fail upload file: " + resp.Status)
	}

	return nil
}

func (sd *SelCDN) getAccessToken() (string, error) {
	if sd.accessToken != "" {
		return sd.accessToken, nil
	}

	req, err := http.NewRequest("GET", sd.AuthURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Auth-User", sd.AuthUser)
	req.Header.Set("X-Auth-Key", sd.AuthKey)

	client := http.Client{Timeout: sd.timeOutRequest * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Log.Error(errors.Wrap(err, "fail close file body"))
		}
	}()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return "", errors.New("fail to get access token")
	}

	key := resp.Header.Get("X-Auth-Token")
	if key == "" {
		return "", errors.New("fail to get access token: key AuthToken is empty")
	}

	sd.accessToken = key
	return key, nil
}

// GetCDNStorage
func GetCDNStorage() ICDNFileStorage {
	return cdnStorage
}
