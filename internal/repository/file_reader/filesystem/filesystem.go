package filesystem

import (
	"fmt"
	"io/ioutil"
	"suse-cli-challenge/internal/repository"

	"context"
)

type FilesystemFileReaderRepository struct {
	repository.FileReaderRepository
}

func (g *FilesystemFileReaderRepository) RetrieveFileContent(ctx context.Context, path string) ([]byte, error) {
	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrReadingFile, err.Error())
	}

	return content, nil
}

func (g *FilesystemFileReaderRepository) RetrieveDirectoryContent(
	ctx context.Context,
	path string,
) (string, error) {
	return path, nil
}
func NewFilesystemFileReaderRepository() repository.FileReaderRepository {
	return &FilesystemFileReaderRepository{}
}
