package utils

import (
	"context"
	"fmt"
	"log"
	"os"
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
		if strings.Contains(*file.Filename, "dashboard") {
			fmt.Println("File:", file.GetFilename())
			return file, nil
		}
	}
	return nil, fmt.Errorf("error: file not found in PR: %d", pull)
}

func GetFileContent(client *github.Client, file *github.CommitFile, owner, repo string) error {
	content, _, _, err := client.Repositories.GetContents(ctx, owner, repo, *file.Filename, nil)
	if err != nil {
		return fmt.Errorf("failed to get file content: %v", err)
	}

	decodedContent, err := content.GetContent()
	if err != nil {
		return fmt.Errorf("failed to decode content: %v", err)
	}

	os.OpenFile("dashboard.yaml", os.O_RDWR|os.O_CREATE, 0755)
	err = os.WriteFile("dashboard.yaml", []byte(decodedContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
