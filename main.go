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
	token = os.Getenv("AUTH_TOKEN")
	ref   = os.Getenv("GITHUB_REF")
	repo  = os.Getenv("GITHUB_REPOSITORY")
	check = os.Getenv("CHECK")
)

func main() {
	flag.Parse()

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

	file, err := u.SelectFile(pull, files)
	if err != nil {
		fmt.Printf("Error selecting file: %v\n", err)
		os.Exit(1)
	}

	l.ExtractJsonFromYamlFile(file)
	if err != nil {
		fmt.Printf("Error extracting json from yaml file: %v\n", err)
		os.Exit(1)
	}

	switch check {
	case "linter":
		fmt.Println("Running linter check")
		results, err := l.LintJsonFile("dashboard.json")
		if err != nil {
			fmt.Printf("Error linting json file: %v\n", err)
			os.Exit(1)
		}
		if results != nil {
			fmt.Println("\nResults:")
			results.ReportByRule()
		}
	case "validator":
		fmt.Println("Running validator check")
		// TODO: Implement validator check here for UID

	case "both":
		fmt.Println("Running both linter and validator checks")
		results, err := l.LintJsonFile("dashboard.json")
		if err != nil {
			fmt.Printf("Error linting json file: %v\n", err)
			os.Exit(1)
		}
		if results != nil {
			fmt.Println("\nResults:")
			results.ReportByRule()
		}
		// TODO: Implement validator check here for UID
	default:
		fmt.Println("Invalid check type")
		os.Exit(1)
	}
}
