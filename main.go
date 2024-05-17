package main

import (
	"fmt"
	"os"

	l "github.com/ministryofjustice/cloud-platform-grafana-linter/linter"
	u "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"
)

type Config struct {
	Owner               string
	Repo                string
	Token               string
	PULL_REQUEST_NUMBER int
}

var (
	githubref = os.Getenv("GITHUB_REF")
	token     = os.Getenv("GITHUB_TOKEN")
	owner     = os.Getenv("GITHUB_OWNER")
	repo      = os.Getenv("GITHUB_REPO")
	// kubeConfigPath = os.Getenv("KUBE_CONFIG_PATH")
)

func main() {
	client, ctx := u.GitHubClient(token)

	pull, err := u.GetPullRequestNumber(client, owner, repo, githubref)

	files, err := u.ListFiles(owner, repo, client, ctx, pull)
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		os.Exit(1)
	}

	file, err := u.SelectFile(pull, files)
	if err != nil {
		fmt.Printf("Error selecting file: %v\n", err)
		os.Exit(1)
	}

	b, results, err := l.ExtractJsonFromYamlFile(file)
	if err != nil {
		fmt.Printf("Error extracting json from yaml file: %v\n", err)
		os.Exit(1)
	}

	if b {
		results.ReportByRule()
	}
}
