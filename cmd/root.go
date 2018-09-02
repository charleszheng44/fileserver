package cmd

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	maxUploadSize = 2 * 1024 * 1024 // 2 Gb
)

var (
	uploadPath string
	addr       string
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fileserver",
		Short: "fileserver start a basic http fileserver",
		Long:  "fileserver support basic upload and download operations",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(addr, uploadPath)
		},
	}
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "127.0.0.1:9344", "network address fileserver will liseten to. (<ip:port>)")
	rootCmd.PersistentFlags().StringVar(&uploadPath, "path", "/data/workspace", "path to the directory that stores uploaded data.")

	return rootCmd
}

func startServer(addr string, uploadPath string) {
	http.HandleFunc("/upload", uploadHandler)
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/download/", http.StripPrefix("/download", fs))
	logrus.Infof("fileserver will start on %s, use /upload to upload files and /download/{filename} to download file", addr)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

// uploadHandler is the handler for uploading file
func uploadHandler(repWriter http.ResponseWriter, req *http.Request) {
	// validate file size
	req.Body = http.MaxBytesReader(repWriter, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		logrus.Errorf("file to big, file to be uploaded can't large than %d byte: %v", maxUploadSize, err)
		renderError(repWriter, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	// parse and validate file and post parameters
	fileName := req.PostFormValue("name")
	file, _, err := req.FormFile("uploadFile")
	if err != nil {
		logrus.Errorf("fail to upload file: %v", err)
		renderError(repWriter, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("fail to upload file: %v", err)
		renderError(repWriter, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	newPath := filepath.Join(uploadPath, fileName)
	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		logrus.Errorf("fail to upload file: %v", err)
		renderError(repWriter, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		logrus.Errorf("fail to upload file: %v", err)
		renderError(repWriter, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	repWriter.Write([]byte("SUCCESS"))
}

// renderError is the helper function for render error to http responseWriter
func renderError(repWriter http.ResponseWriter, message string, statusCode int) {
	repWriter.WriteHeader(http.StatusBadRequest)
	repWriter.Write([]byte(message))
}
