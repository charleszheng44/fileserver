package fileserver

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	maxUploadSize = 2 * 1024 * 1024 // 2 Gb
)

type FileServer struct {
	UploadHandler func(http.ResponseWriter, *http.Request)
	UploadPath    string
	DownloadPath  string
	Address       string
}

// NewFileServer returns a new FileServer pointer
func NewFileServer(addr, uploadPath, downloadPath string) *FileServer {
	fs := &FileServer{
		UploadPath:   uploadPath,
		DownloadPath: downloadPath,
		Address:      addr,
	}
	fs.RegisterUploadHandler()
	return fs
}

// StartServer starts a fileserver instance
func (fs *FileServer) StartServer(addr, uploadPath string) {
	http.HandleFunc("/upload", fs.UploadHandler)
	fsInst := http.FileServer(http.Dir(fs.UploadPath))
	http.Handle("/download/", http.StripPrefix("/download", fsInst))
	logrus.Infof(`fileserver will start on %s, use /upload to upload 
files and /download/{filename} to download file`, fs.Address)
	logrus.Fatal(http.ListenAndServe(fs.Address, nil))
}

func (fs *FileServer) RegisterUploadHandler() {
	fs.UploadHandler = func(repWriter http.ResponseWriter, req *http.Request) {
		// validate file size
		req.Body = http.MaxBytesReader(repWriter, req.Body, maxUploadSize)
		if err := req.ParseMultipartForm(maxUploadSize); err != nil {
			logrus.Errorf(`file to big, file to be uploaded can't 
large than %d byte: %v`, maxUploadSize, err)
			renderError(repWriter, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}

		// parse and validate file and post parameters
		fileName := req.PostFormValue("name")
		fileModeStr, parseErr :=
			strconv.ParseUint(req.PostFormValue("mode"), 10, 32)
		if parseErr != nil {
			logrus.Errorf("fail to get upload file mode: %v", parseErr)
			renderError(repWriter, "INVALID_FILE_MODE", http.StatusBadRequest)
			return
		}
		fileMode := uint32(fileModeStr)

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

		newPath := filepath.Join(fs.UploadPath, fileName)
		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			logrus.Errorf("fail to upload file: %v", err)
			renderError(repWriter, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			logrus.Errorf("fail to upload file: %v", err)
			renderError(repWriter, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}

		if err = os.Chmod(newPath, os.FileMode(fileMode)); err != nil {
			logrus.Errorf("fail to change file mode: %v", err)
			renderError(repWriter, "CANT_CHANGE_FILE_MODE", http.StatusInternalServerError)
			return
		}

		repWriter.Write([]byte("SUCCESS"))
	}
}

// renderError is the helper function for render error to http responseWriter
func renderError(repWriter http.ResponseWriter, message string, statusCode int) {
	repWriter.WriteHeader(http.StatusBadRequest)
	repWriter.Write([]byte(message))
}
