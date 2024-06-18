package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	token = os.Getenv("AUTH_TOKEN")
	ref   = os.Getenv("GITHUB_REF")
	repo  = os.Getenv("GITHUB_REPOSITORY")
)

func main() {
	githubrefS := strings.Split(ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	client := u.GitHubClient(token)

	files, _, err := u.GetPullRequestFiles(client, owner, repoName, pull)
	if err != nil {
		fmt.Printf("Error fetching files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("GetPullRequestFiles: Done")

	file, err := u.SelectFile(pull, files)
	if err != nil {
		fmt.Printf("Error selecting file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SelectFile: Done")

	b, results, err := l.ExtractJsonFromYamlFile(file)
	if err != nil {
		fmt.Printf("Error extracting json from yaml file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ExtractJsonFromYamlFile: Done")

	if b {
		results.ReportByRule()
	}
}
