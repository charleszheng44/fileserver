package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type Client struct {
	client    *http.Client
	remoteURL string
}

func NewClient(remoteURL string) *Client {
	return &Client{
		client:    &http.Client{},
		remoteURL: remoteURL,
	}
}

// Upload uploads the given file to the fileserver
func (c *Client) UpLoad(filePath string) error {
	remoteFileName := filepath.Base(filePath)
	fileDescriptor, err := os.Open(filePath)
	if err != nil {
		logrus.Errorf("fail to open file %s: %v", filePath, err)
		return err
	}
	uploadFileMode, err := getFileMode(filePath)
	if err != nil {
		logrus.Errorf("fail to get file state: %v", err)
		return err
	}

	// prepare the request body
	values := map[string]io.Reader{
		"uploadFile": fileDescriptor,
		"name":       strings.NewReader(remoteFileName),
		"mode":       strings.NewReader(fmt.Sprint(uint32(uploadFileMode))), // the file mode is uint32
	}

	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)

	// generate the request body
	for key, val := range values {
		var fieldWriter io.Writer
		if r, ok := val.(io.Closer); ok {
			defer r.Close()
		}

		if r, ok := val.(*os.File); ok {
			// add file content
			if fieldWriter, err = multipartWriter.CreateFormFile(key, r.Name()); err != nil {
				return err
			}
		} else {
			// add other field
			if fieldWriter, err = multipartWriter.CreateFormField(key); err != nil {
				return err
			}
		}
		// copy the file/field contents to multipartwriter
		if _, err = io.Copy(fieldWriter, val); err != nil {
			return err
		}
	}

	// don't forget to close the multipartWriter
	if closeErr := multipartWriter.Close(); closeErr != nil {
		logrus.Warnf("fail to close multipart writer: %v", closeErr)
	}

	// generate new request
	uploadURL := fmt.Sprintf("%s/upload", c.remoteURL)
	postReq, err := http.NewRequest("POST", uploadURL, &buf)
	if err != nil {
		return err
	}

	// set the content type
	postReq.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// send the request
	postRep, err := c.client.Do(postReq)
	if err != nil {
		return err
	}
	defer postRep.Body.Close()

	// check response
	if postRep.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", postRep.Status)
		return err
	}

	return nil
}

// Download downloads the desired file from fileserver
func (c *Client) Download(filename string) error {

	fileDescriptor, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer fileDescriptor.Close()

	downloadUrl := fmt.Sprintf("%s/download/%s", c.remoteURL, filename)

	getRep, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer getRep.Body.Close()

	_, err = io.Copy(fileDescriptor, getRep.Body)
	if err != nil {
		return err
	}

	return nil
}

func getFileMode(filename string) (os.FileMode, error) {
	fileInfo, statErr := os.Stat(filename)
	if statErr != nil {
		return os.FileMode(0), statErr
	}
	return fileInfo.Mode(), nil
}
