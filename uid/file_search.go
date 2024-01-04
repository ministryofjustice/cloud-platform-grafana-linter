package uid

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
)

var (
	owner = "jackstockley89"
	repo  = "cloud-platform-environments"
)

// ListFiles retrieves a list of commit files for each pull request in a GitHub repository.
// It takes a GitHub client and a context as input parameters.
// It returns a slice of commit files, and an error if any.
func ListFiles(client *github.Client, ctx context.Context, pull int) ([]*github.CommitFile, error) {
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pull, nil)
	if err != nil {
		error := fmt.Errorf("error: listing files: %v", err)
		return nil, error
	}
	return files, nil
}

func SearchCode(pull int, files []*github.CommitFile) (string, error) {
	// Iterate over each file
	for _, file := range files {
		// pass file.Patch to a function to parse the patch
		// and return the UID

		filePatch := *file.Patch

		// Split the patch into lines
		lines := strings.Split(filePatch, "\n")

		// Iterate over each line
		for _, line := range lines {
			// Check if the line contains the search string
			if strings.Contains(line, "uid") {
				fmt.Println("Found in PR:", pull)
				fmt.Println("File:", *file.Filename)
				word := strings.Fields(line)
				// trim the trailing comma
				uid := strings.TrimSuffix(word[2], ",")
				// trim the leading double quote
				uid = strings.TrimPrefix(uid, "\"")
				// trim the trailing double quote
				uid = strings.TrimSuffix(uid, "\"")

				fmt.Println("UID:", uid)

				return uid, nil
			}
		}
	}
	return "", fmt.Errorf("error: UID not found in PRs")
}
