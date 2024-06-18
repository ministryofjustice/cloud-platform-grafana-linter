package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v62/github"
)

var (
	ctx = context.Background()
)

func GitHubClient(token string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)

	_, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil
	}

	// Rate.Limit should most likely be 5000 when authorized.
	log.Printf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		log.Printf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	return client
}

// ListFiles retrieves a list of commit files for each pull request in a GitHub repository.
// It takes a GitHub client and a context as input parameters.
// It returns a slice of commit files, and an error if any.
func GetPullRequestFiles(client *github.Client, o, r string, n int) ([]*github.CommitFile, *github.Response, error) {
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
			return file, nil
		}
	}
	return nil, fmt.Errorf("error: file not found in PR: %d", pull)
}
