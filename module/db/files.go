package db

import (
	"github.com/lexbond13/api_core/module/db/structure"
	"time"
)

type IFilesRepository interface {
	Create(file *structure.File) error
	Delete(ID int64) error
}

// NewFileRepository
func NewFileRepository() IFilesRepository {
	return &FileRepository{}
}

type FileRepository struct {
}

// Create
func (f *FileRepository) Create(file *structure.File) error {
	file.CreatedAt = time.Now()
	_, err:= connection.Model(file).Insert()
	if err != nil {
		return err
	}
	return nil
}

// Delete
func (f *FileRepository) Delete(ID int64) error {
	fileStruct := &structure.File{ID:ID}
	_, err := connection.Model(fileStruct).WherePK().Delete()

	return err
}
