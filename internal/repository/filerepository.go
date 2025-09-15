package repository

import (
	"github.com/ruslannnnnnnnn/test-file-storage/internal/model"
	"gorm.io/gorm"
)

type IFileRepository interface {
	AutoMigrate() error
	ListFiles() ([]model.File, error)
}

type FileRepository struct {
	dbConnection *gorm.DB
}

func NewFileRepository(dbConnection *gorm.DB) IFileRepository {
	return &FileRepository{dbConnection: dbConnection}
}

func (f FileRepository) AutoMigrate() error {
	err := f.dbConnection.AutoMigrate(model.File{})
	if err != nil {
		return err
	}

	return nil
}

func (f FileRepository) ListFiles() ([]model.File, error) {
	var files []model.File

	result := f.dbConnection.Find(&files, "true")
	if result.Error != nil {
		return nil, result.Error
	}

	return files, nil
}
