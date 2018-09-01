A simple http fileserver implement in go

#### HOW TO RUN 

```bash
go run fileserver.go --addr 127.0.0.1:9344
```
### SAMPLE HTTP REQUEST

1. Upload a file

```bash
curl -X POST -F "uploadFile=@/path/to/the/file" -F "filename=<name_of_file_on_fileserver>" http://127.0.0.1:9344
```

2. Download a file

```bash
curl -X GET http://127.0.0.1:9344/download/<filename>
```
