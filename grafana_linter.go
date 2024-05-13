package main

import (
	"fmt"
	"os"
	"strconv"

	uc "github.com/ministryofjustice/cloud-platform-grafana-linter/cluster"
	l "github.com/ministryofjustice/cloud-platform-grafana-linter/linter"
	uid "github.com/ministryofjustice/cloud-platform-grafana-linter/uid"
	u "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"
)

var (
	p = os.Getenv("PULL_REQUEST_NUMBER")
)

func main() {
	client, ctx := u.GitHubClient()

	// convert pull output to int value
	pull, err := strconv.Atoi(p)
	if err != nil {
		fmt.Printf("Error converting pull request number to int: %v\n", err)
		os.Exit(1)
	}

	files, err := uid.ListFiles(client, ctx, pull)
	if err != nil {
		fmt.Printf("Error listing files: %v\n", err)
		os.Exit(1)
	}

	uid, err := uid.SearchCode(pull, files)
	if err != nil {
		fmt.Printf("Error searching code: %v\n", err)
		os.Exit(1)
	}

	clientset := uc.ClientSet()
	configMaps, err := uc.SearchNamespacesForConfigMaps(clientset)
	if err != nil {
		fmt.Printf("Error getting configmaps: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %v configmaps\n", len(configMaps))
	for _, configMap := range configMaps {
		d := uc.SearchConfigMapsData(clientset, configMap, uid)
		if d != "" {
			fmt.Printf("%v", d)
		}
	}

	l.ExtractJsonFromPullRequestFile(files)
}
