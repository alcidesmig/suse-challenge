package filesystem

import (
	"fmt"
	"io/ioutil"
	"suse-cli-challenge/internal/repository"

	"context"
)

// FilesystemFileReaderRepository represents a file reader repository implementation for local filesystem.
type FilesystemFileReaderRepository struct {
	repository.FileReaderRepository
}

// RetrieveFileContent retrieves the content of a file from the filesystem.
func (g *FilesystemFileReaderRepository) RetrieveFileContent(ctx context.Context, path string) ([]byte, error) {
	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", repository.ErrReadingFile, err.Error())
	}

	return content, nil
}

// RetrieveDirectoryContent retrieves the path the content of a directory from the filesystem.
func (g *FilesystemFileReaderRepository) RetrieveDirectoryContent(
	ctx context.Context,
	path string,
) (string, error) {
	return path, nil
}
func NewFilesystemFileReaderRepository() repository.FileReaderRepository {
	return &FilesystemFileReaderRepository{}
}
