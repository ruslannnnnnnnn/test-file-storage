package service

import (
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
)

type IFileService interface {
	Upload(common.UploadRequest) (common.UploadResponse, error)
	Download(common.DownloadRequest) (common.DownloadResponse, error)
	ListFiles(common.ListFilesRequest) (common.ListFilesResponse, error)
}

type FileService struct {
}

func (f FileService) Upload(request common.UploadRequest) (common.UploadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileService) Download(request common.DownloadRequest) (common.DownloadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileService) ListFiles(request common.ListFilesRequest) (common.ListFilesResponse, error) {
	//TODO implement me
	panic("implement me")
}
