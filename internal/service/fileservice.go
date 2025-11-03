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

func (f *FileService) Upload(ctx context.Context, filename string, fileIdCh chan string, r io.Reader) error {
	fileId, err := f.fileRepository.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка при запросе на создание записи о файле %w", err)
	}

	select {
	case fileIdCh <- fileId:
	case <-ctx.Done():
		return ctx.Err()
	}

	err = filesystem.WriteStream(fileId, r)
	if err != nil {
		return fmt.Errorf("ошибка при записи файла на диск %w", err)
	}

	return nil
}

func (f *FileService) Download(ctx context.Context, id string) (string, io.Reader, error) {
	err := uuid.Validate(id)
	if err != nil {
		return "", nil, InvalidUUidError
	}

	file, err := f.fileRepository.GetById(id)
	if err != nil {
		return "", nil, FileNotFoundError
	}

	r, w := io.Pipe()

	go func() {
		defer w.Close()
		err = filesystem.Read(file.Name, w)
		if err != nil {
			w.CloseWithError(err)
		}
	}()

	return file.Name, r, nil

}

func (f *FileService) ListFiles(ctx context.Context) (common.ListFilesResponse, error) {
	result, err := f.fileRepository.ListFiles()
	if err != nil {
		return common.ListFilesResponse{}, fmt.Errorf("ошибка при запросе списка файлов %w", err)
	}

	return common.ListFilesResponse{Files: result}, nil
}
