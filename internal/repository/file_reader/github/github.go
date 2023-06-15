package github

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"suse-cli-challenge/internal/repository"

	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v37/github"
)

// GithubFileReaderRepository represents a file reader repository implementation for GitHub.
type GithubFileReaderRepository struct {
	GithubClient *github.Client
	repository.FileReaderRepository
}

// RetrieveFileContent retrieves the content of a file from one GitHub repository.
func (g *GithubFileReaderRepository) RetrieveFileContent(ctx context.Context, url string) ([]byte, error) {
	repoData, err := g.parseURL(url)
	if err != nil {
		return nil, errParsingGithubPath(repository.ErrInvalidRepoURL)
	}
	content, _, _, err := g.GithubClient.Repositories.GetContents(ctx, repoData.Owner, repoData.RepositoryName, repoData.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrLoadingChart, err.Error())
	}

	if !strings.EqualFold(content.GetType(), "file") {
		return nil, fmt.Errorf("%w: %s", repository.ErrInvalidChartDir, fmt.Sprintf("error in Chart.yaml format. File expected, got %s", content.GetType()))
	}

	data, err := content.GetContent()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", repository.ErrLoadingChart, err.Error())
	}
	return []byte(data), nil
}

// RetrieveDirectoryContent retrieves the content of a directory from a specific GitHub repository.
// It clones the repository to a temporary location, removes any existing directory with the same name,
// and then returns the path of the cloned directory where the desired content is located.
func (g *GithubFileReaderRepository) RetrieveDirectoryContent(
	ctx context.Context,
	repoURL string,
) (string, error) {
	fmt.Println("Cloning repository for packaging chart...")

	// Parse the GitHub repository URL
	repoData, err := g.parseURL(repoURL)
	if err != nil {
		return "", err
	}

	// Prepare the location for cloning the repository
	chartLocation := fmt.Sprintf("%s/%s", os.TempDir(), repoData.RepositoryName)

	// Remove any existing directory with the same name
	os.RemoveAll(chartLocation)

	// Clone the repository to the specified location
	repoSanitizedURL := fmt.Sprintf("https://github.com/%s/%s", repoData.Owner, repoData.RepositoryName)
	_, err = git.PlainClone(chartLocation, false, &git.CloneOptions{
		URL:           repoSanitizedURL,
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(repoData.Branch),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrCloningRepo, err.Error())
	}

	// Return the path of the cloned directory where the desired content is located
	return fmt.Sprintf("%s/%s", chartLocation, repoData.Path), nil
}

// NewGithubFileReaderRepository creates a new instance of FileReaderRepository for GithubFileReaderRepository implementation.
func NewGithubFileReaderRepository() repository.FileReaderRepository {
	httpClient := http.Client{}
	return &GithubFileReaderRepository{GithubClient: github.NewClient(&httpClient)}
}
