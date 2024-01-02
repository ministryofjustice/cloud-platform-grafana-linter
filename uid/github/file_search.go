package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

var (
	owner = "jackstockley89"
	repo  = "cloud-platform-environments"
	pulls = []int{52}
	token = os.Getenv("GITHUB_TOKEN")
)

func listFiles(client *github.Client, ctx context.Context) (int, []*github.CommitFile, error) {
	// Iterate over each pull request
	for _, pull := range pulls {
		// List files for each pull request
		files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, pull, nil)
		if err != nil {
			error := fmt.Errorf("error: listing files: %v", err)
			return 0, nil, error
		}
		return pull, files, nil
	}
	return 0, nil, fmt.Errorf("error: listing files: %v", "no pull requests found")
}

func Client() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx
}

func SearchCode(client *github.Client, ctx context.Context) (string, error) {
	pull, files, err := listFiles(client, ctx)
	if err != nil {
		error := fmt.Errorf("error: listing files: %v", err)
		return "", error
	}

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

func ExtractJsonFromPullRequestFile(client *github.Client, ctx context.Context) (string, error) {
	_, files, err := listFiles(client, ctx)
	if err != nil {
		error := fmt.Errorf("error: listing files: %v", err)
		return "", error
	}
	// Iterate over each file
	for _, file := range files {
		// pass file.Patch into a local json file

		filePatch := *file.Patch

		jsonFile, err := os.Create("dashboard.json")
		if err != nil {
			error := fmt.Errorf("error: creating json file: %v", err)
			return "", error
		}

		defer jsonFile.Close()

		_, err = jsonFile.WriteString(filePatch)
		if err != nil {
			error := fmt.Errorf("error: writing to json file: %v", err)
			return "", error
		}

	}
	return "dashboard.json", nil
}
