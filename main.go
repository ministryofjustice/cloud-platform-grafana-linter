package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	l "github.com/ministryofjustice/cloud-platform-grafana-linter/linter"
	u "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"
)

var (
	token = flag.String(os.Getenv("AUTH_TOKEN"), "token", "GitHub token used for authentication")
	ref   = flag.String(os.Getenv("GITHUB_REF"), "ref", "GitHub pull request ref (e.g. refs/pull/1/head)")
	repo  = flag.String(os.Getenv("GITHUB_REPOSITORY"), "repo", "GitHub repository (e.g. owner/repository)")
	check = flag.String(os.Getenv("CHECK"), "check", "Check for selecting linter or validator")
)

func main() {
	flag.Parse()

	githubrefS := strings.Split(*ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(*repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	client := u.GitHubClient(*token)

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

	if *check == "linter" {
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

	if *check == "validator" {
		// TODO: Implement validator check here for UID
		fmt.Println("Validator check not implemented yet")
	}
}
