package fileserver

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
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

func uploadHandler(repWriter http.ResponseWriter, req *http.Request) {
}
