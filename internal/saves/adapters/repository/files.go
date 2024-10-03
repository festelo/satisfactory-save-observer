package repository

import (
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/festelo/satisfactory-save-observer/internal/saves/domain"
)

type FilesSaveRepository struct {
	directory string
}

func NewFilesSaveRepository(directory string) *FilesSaveRepository {
	return &FilesSaveRepository{directory}
}

func (r FilesSaveRepository) ListSaves() ([]*domain.Save, error) {
	files, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, err
	}

	saves := []*domain.Save{}

	for _, file := range files {
		if file.Type().IsRegular() {
			name := file.Name()
			ext := filepath.Ext(name)
			if ext == ".sav" {
				fileInfo, err := os.Stat(path.Join(r.directory, name))
				if err != nil {
					return nil, err
				}

				saves = append(saves, &domain.Save{
					Name:      strings.TrimSuffix(name, ext),
					CreatedAt: fileInfo.ModTime(),
				})
			}

		}
	}

	return saves, nil
}

func (r FilesSaveRepository) CopySave(name string, writer io.Writer) error {
	savePath := path.Join(r.directory, name+".sav")

	file, err := os.Open(savePath)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			slog.Warn("Error closing %v: %v", savePath, err)
		}
	}()

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}
