package github

import (
	"fmt"
	"net/url"
	"strings"
	"suse-cli-challenge/internal/models"
	"suse-cli-challenge/internal/repository"
)

func errParsingGithubPath(err error) error {
	return fmt.Errorf("%w: error parsing github path. Expected: https://github.com/<owner>/<repo>/tree/<branch>/<path>", err)
}

// parseURL parses the GitHub repository URL and extracts owner, repository name, path, and branch.
func (g *GithubFileReaderRepository) parseURL(repoURL string) (*models.GithubMetadata, error) {
	// Parse the URL
	u, err := url.Parse(repoURL)
	if err != nil {
		return nil, err
	}

	// Check if the host is GitHub
	if u.Host != "github.com" {
		return nil, errParsingGithubPath(repository.ErrInvalidRepoURL)
	}

	// Remove the leading slash from the path
	path := strings.TrimPrefix(u.Path, "/")

	// Split the path into parts
	parts := strings.Split(path, "/")

	// Extract owner, repository name, and path
	if len(parts) < 2 {
		return nil, errParsingGithubPath(repository.ErrInvalidRepoURL)
	}
	owner := parts[0]
	repo := parts[1]
	branch := parts[3]

	// Starting in 4 remove the /tree/<branch> section
	path = strings.Join(parts[4:], "/")

	return &models.GithubMetadata{Owner: owner, RepositoryName: repo, Path: path, Branch: branch}, nil
}
