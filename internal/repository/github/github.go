package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"suse-cli-challenge/internal/repository"

	"context"

	"suse-cli-challenge/internal/models"

	"github.com/google/go-github/v37/github"
)

type GithubRepository struct {
	GithubClient *github.Client
	repository.GithubRepository
}

func (g *GithubRepository) ParseURL(repoURL string) (*models.GithubMetadata, error) {
	// Parse the URL
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}

	// Check if the host is GitHub
	if u.Host != "github.com" {
		return nil, fmt.Errorf("invalid GitHub URL: %s", repoURL)
	}

	// Remove the leading slash from the path
	path := strings.TrimPrefix(u.Path, "/")

	// Split the path into parts
	parts := strings.Split(path, "/")

	// Extract owner, repository name, and path
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub URL: %s", repoURL)
	}
	owner := parts[0]
	repo := parts[1]
	// Starting in 4 remove the /tree/<branch> section
	path = strings.Join(parts[4:], "/")

	return &models.GithubMetadata{Owner: owner, RepositoryName: repo, Path: path}, nil
}

func (g *GithubRepository) RetrieveFileContent(ctx context.Context, url string) ([]byte, error) {
	repoData, err := g.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>")
	}
	content, _, _, err := g.GithubClient.Repositories.GetContents(ctx, repoData.Owner, repoData.RepositoryName, repoData.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving %s: %s", repoData.Path, err.Error())
	}

	if !strings.EqualFold(content.GetType(), "file") {
		return nil, fmt.Errorf("error in Chart.yaml format. File expected, got %s", content.GetType())
	}

	data, err := content.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %s", err.Error())
	}
	return []byte(data), nil
}

func NewGithubRepositoy() *GithubRepository {
	httpClient := http.Client{}
	return &GithubRepository{GithubClient: github.NewClient(&httpClient)}
}
