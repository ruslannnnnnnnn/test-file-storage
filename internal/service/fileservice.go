package service

import (
	"github.com/google/uuid"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/filesystem"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/repository"
)

type IFileService interface {
	Upload(common.UploadRequest) (string, error)
	Download(common.DownloadRequest) (common.DownloadResponse, error)
	ListFiles() (common.ListFilesResponse, error)
}

type FileService struct {
	fileRepository repository.IFileRepository
}

func NewFileService(fileRepository repository.IFileRepository) *FileService {
	return &FileService{fileRepository: fileRepository}
}

func (f FileService) Upload(request common.UploadRequest) (fileId string, err error) {

	fileId, err = f.fileRepository.Create(request.FileName)
	if err != nil {
		return "", err
	}

	// у сохранённого файла uuid в названии чтобы можно было иметь несколько файлов с одинаковым названием
	err = filesystem.Write(fileId, request.FileContent)
	if err != nil {
		return "", err
	}

	return
}

func (f FileService) Download(request common.DownloadRequest) (common.DownloadResponse, error) {

	err := uuid.Validate(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, InvalidUUidError{}
	}

	file, err := f.fileRepository.GetById(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, FileNotFoundError{}
	}

	fileContent, err := filesystem.Read(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, err
	}

	return common.DownloadResponse{FileName: file.Name, FileContent: fileContent}, nil
}

func (f FileService) ListFiles() (common.ListFilesResponse, error) {
	result, err := f.fileRepository.ListFiles()
	if err != nil {
		return common.ListFilesResponse{}, err
	}

	return common.ListFilesResponse{Files: result}, nil
}
