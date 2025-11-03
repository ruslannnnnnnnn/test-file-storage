package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/common"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/filesystem"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/repository"
)

var (
	InvalidUUidError  = errors.New("invalid uuid")
	FileNotFoundError = errors.New("file not found")
)

type IFileService interface {
	Upload(ctx context.Context, filename string, fileIdCh chan string, r io.Reader) error
	Download(ctx context.Context, id string) (string, io.Reader, error)
	ListFiles(ctx context.Context) (common.ListFilesResponse, error)
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
		return "", fmt.Errorf("ошибка при запросе на создание записи о файле %w", err)
	}

	// у сохранённого файла uuid в названии чтобы можно было иметь несколько файлов с одинаковым названием
	err = filesystem.Write(fileId, request.FileContent)
	if err != nil {
		return "", fmt.Errorf("ошибка при при попытке записи файла в файловой системе %w", err)
	}

	return
}

func (f FileService) Download(request common.DownloadRequest) (common.DownloadResponse, error) {

	err := uuid.Validate(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, InvalidUUidError
	}

	file, err := f.fileRepository.GetById(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, FileNotFoundError
	}

	fileContent, err := filesystem.Read(request.FileId)
	if err != nil {
		return common.DownloadResponse{}, fmt.Errorf("ошибка при попытке чтения файла из файловой системы %w", err)
	}

	return common.DownloadResponse{FileName: file.Name, FileContent: fileContent}, nil
}

func (f FileService) ListFiles() (common.ListFilesResponse, error) {
	result, err := f.fileRepository.ListFiles()
	if err != nil {
		return common.ListFilesResponse{}, fmt.Errorf("ошибка при запросе списка файлов %w", err)
	}

	return common.ListFilesResponse{Files: result}, nil
}
