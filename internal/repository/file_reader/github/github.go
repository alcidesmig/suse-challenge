package github

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"suse-cli-challenge/internal/repository"

	"context"

	"suse-cli-challenge/internal/models"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-github/v37/github"
)

type GithubFileReaderRepository struct {
	GithubClient *github.Client
	repository.FileReaderRepository
}

func (g *GithubFileReaderRepository) parseURL(repoURL string) (*models.GithubMetadata, error) {
	// Parse the URL
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}

	// Check if the host is GitHub
	if u.Host != "github.com" {
		return nil, fmt.Errorf("%w: error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>", repository.ErrInvalidRepoURL)
	}

	// Remove the leading slash from the path
	path := strings.TrimPrefix(u.Path, "/")

	// Split the path into parts
	parts := strings.Split(path, "/")

	// Extract owner, repository name, and path
	if len(parts) < 2 {
		return nil, fmt.Errorf("%w: error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>", repository.ErrInvalidRepoURL)
	}
	owner := parts[0]
	repo := parts[1]
	branch := parts[3]
	// Starting in 4 remove the /tree/<branch> section
	path = strings.Join(parts[4:], "/")

	return &models.GithubMetadata{Owner: owner, RepositoryName: repo, Path: path, Branch: branch}, nil
}

func (g *GithubFileReaderRepository) RetrieveFileContent(ctx context.Context, url string) ([]byte, error) {
	repoData, err := g.parseURL(url)
	if err != nil {
		return nil, fmt.Errorf("%w: error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>", repository.ErrInvalidRepoURL)
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

func (g *GithubFileReaderRepository) RetrieveDirectoryContent(
	ctx context.Context,
	repoURL string,
) (string, error) {
	fmt.Println("Cloning repository for packaging chart...")
	repoData, err := g.parseURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("%w: error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>", repository.ErrInvalidRepoURL)
	}

	chartLocation := fmt.Sprintf("%s/%s", os.TempDir(), repoData.RepositoryName)
	os.RemoveAll(chartLocation)

	_, err = git.PlainClone(chartLocation, false, &git.CloneOptions{
		URL:           fmt.Sprintf("https://github.com/%s/%s", repoData.Owner, repoData.RepositoryName),
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(repoData.Branch),
	})
	if err != nil {
		return "", fmt.Errorf("%w: %s", repository.ErrCloningRepo, err.Error())
	}

	return fmt.Sprintf("%s/%s", chartLocation, repoData.Path), nil
}

func NewGithubFileReaderRepository() repository.FileReaderRepository {
	httpClient := http.Client{}
	return &GithubFileReaderRepository{GithubClient: github.NewClient(&httpClient)}
}
