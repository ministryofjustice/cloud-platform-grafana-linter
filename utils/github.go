package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

var (
	ctx = context.Background()
)

func GitHubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

// ListFiles retrieves a list of commit files for each pull request in a GitHub repository.
// It takes a GitHub client and a context as input parameters.
// It returns a slice of commit files, and an error if any.
func GetPullRequestFiles(token, o, r string, n int) ([]*github.CommitFile, *github.Response, error) {
	client := GitHubClient(token)
	files, resp, err := client.PullRequests.ListFiles(ctx, o, r, n, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching files: %w", err)
	}

	return files, resp, err
}

func SelectFile(pull int, files []*github.CommitFile) (*github.CommitFile, error) {
	for _, file := range files {
		fmt.Println("File:", *file.Filename)
		if strings.Contains(*file.Filename, "dashboard") {
			return nil, nil
		}
	}
	return nil, fmt.Errorf("error: file not found in PR: %d", pull)
}
