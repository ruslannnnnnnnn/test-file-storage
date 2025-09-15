package service

import (
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/repository"
)

type IFileService interface {
	Upload(common.UploadRequest) (common.UploadResponse, error)
	Download(common.DownloadRequest) (common.DownloadResponse, error)
	ListFiles() (common.ListFilesResponse, error)
}

type FileService struct {
	fileRepository repository.IFileRepository
}

func NewFileService(fileRepository repository.IFileRepository) *FileService {
	return &FileService{fileRepository: fileRepository}
}

func (f FileService) Upload(request common.UploadRequest) (common.UploadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileService) Download(request common.DownloadRequest) (common.DownloadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileService) ListFiles() (common.ListFilesResponse, error) {
	result, err := f.fileRepository.ListFiles()
	if err != nil {
		return common.ListFilesResponse{}, err
	}

	return common.ListFilesResponse{Files: result}, nil
}
