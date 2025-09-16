package common

import "github.com/ruslannnnnnnnn/test-file-storage/internal/model"

type DatabaseConfig struct {
	Host   string
	User   string
	DbName string
	Port   int
}

type UploadRequest struct {
	FileName    string
	FileContent []byte
}
type UploadResponse struct {
	FileId string
}

type DownloadRequest struct{}
type DownloadResponse struct{}

type ListFilesResponse struct {
	Files []model.File
}
