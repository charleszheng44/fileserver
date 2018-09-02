package main

import (
	"os"

	"github.com/charleszheng44/fileserver/pkg/client"
	"github.com/sirupsen/logrus"
)

const (
	remoteURL   = "http://127.0.0.1:9344"
	fileName    = "test_upload.txt"
	fileContent = "this file will be uploaded to file server"
)

func main() {
	// create test file
	fd, openFileErr := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0755)
	if openFileErr != nil {
		logrus.Fatalf("fail to open file: %v", openFileErr)
	}
	// write contents to test file
	if _, writeErr := fd.WriteString(fileContent); writeErr != nil {
		logrus.Fatalf("fail to wirte file: %v", writeErr)
	}
	fd.Close()

	fsCli := client.NewClient(remoteURL)
	if uploadErr := fsCli.UpLoad(fileName); uploadErr != nil {
		logrus.Fatalf("fail to upload file: %v", uploadErr)
	}
	logrus.Infof("successfully upload file %s", fileName)

	// remove uploaded file
	if rmErr := os.Remove(fileName); rmErr != nil {
		logrus.Fatalf("fail to remove file: %v", rmErr)
	}

	if downloadErr := fsCli.Download(fileName); downloadErr != nil {
		logrus.Fatalf("fail to download file: %v", downloadErr)
	}

	logrus.Infof("successfully download file %s", fileName)

	if _, fileStatErr := os.Stat(fileName); os.IsNotExist(fileStatErr) {
		logrus.Fatalf("File %s now exist", fileName)
	}
}
