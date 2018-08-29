package fileserver

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	maxUploadSize  = 2 * 1024 * 1024 // 2 Gb
	port           = 9344
	fileserverPath = "/data/workspace"
)

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fs := http.FileServer(http.Dir(fileserverPath))
	http.Handle("/download/", http.StripPrefix("/download", fs))

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	logrus.Infof("file server will start on %s, use /upload to upload files and /download/{filename} to download file", addr)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}

// uploadHandler is the handler for uploading file
func uploadHandler(repWriter http.ResponseWriter, req *http.Request) {
	// validate file size
	req.Body = http.MaxBytesReader(repWriter, req.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	// parse and validate file and post parameters
	fileName := req.PostFormValue("name")
	file, _, err := req.FormFile("uploadFile")
	if err != nil {
		renderError(repWriter, "INVALID_FILE", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		renderError(repWriter, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	if err != nil {
		renderError(repWriter, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		return
	}
	newPath := filepath.Join(uploadPath, fileName)
	fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

	// write file
	newFile, err := os.Create(newPath)
	if err != nil {
		renderError(repWriter, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	defer newFile.Close() // idempotent, okay to call twice
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("SUCCESS"))
}

// renderError is the helper function for render error to http responseWriter
func renderError(repWriter http.ResponseWriter, message string, statusCode int) {
	repWriter.WriteHeader(http.StatusBadRequest)
	repWriter.Write([]byte(message))
}
