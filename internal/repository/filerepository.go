package repository

import (
	"github.com/google/uuid"
	"github.com/ruslannnnnnnnn/test-file-storage/internal/model"
	"gorm.io/gorm"
)

type IFileRepository interface {
	AutoMigrate() error
	ListFiles() ([]model.File, error)
	Create(name string) (id string, err error)
	GetById(id string) (model.File, error)
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

func (f FileRepository) Create(name string) (id string, err error) {

	fileId := uuid.New().String()

	file := model.File{
		Id:   fileId,
		Name: name,
	}

	result := f.dbConnection.Create(&file)

	if result.Error != nil {
		return "", result.Error
	}

	return fileId, nil
}

func (f FileRepository) GetById(id string) (model.File, error) {
	var file model.File
	result := f.dbConnection.First(&file, "id = ?", id)

	if result.Error != nil {
		return model.File{}, result.Error
	}

	return file, nil
}
