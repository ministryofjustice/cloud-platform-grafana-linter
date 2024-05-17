package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

func GitHubClient(token string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx
}

// ListFiles retrieves a list of commit files for each pull request in a GitHub repository.
// It takes a GitHub client and a context as input parameters.
// It returns a slice of commit files, and an error if any.
func ListFiles(owner, repo string, client *github.Client, ctx context.Context, pull int) ([]*github.CommitFile, error) {
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pull, nil)
	if err != nil {
		error := fmt.Errorf("error: listing files: %v", err)
		return nil, error
	}
	return files, nil
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
