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

type DownloadRequest struct {
	FileId string
}
type DownloadResponse struct {
	FileName    string
	FileContent []byte
}

type ListFilesResponse struct {
	Files []model.File
}
