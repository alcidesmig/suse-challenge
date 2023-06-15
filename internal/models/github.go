package models

// GithubMetadata represents some metadata of a GitHub repository.
type GithubMetadata struct {
	Owner          string
	RepositoryName string
	Path           string
	Branch         string
}
