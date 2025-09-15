package common

type DatabaseConfig struct {
	Host   string
	User   string
	DbName string
	Port   int
}

type UploadRequest struct{}
type UploadResponse struct{}

type DownloadRequest struct{}
type DownloadResponse struct{}

type ListFilesRequest struct{}
type ListFilesResponse struct{}
