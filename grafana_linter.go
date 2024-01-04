package main

import (
	"fmt"
	"os"

	uc "github.com/ministryofjustice/cloud-platform-grafana-linter/cluster"
	l "github.com/ministryofjustice/cloud-platform-grafana-linter/linter"
	uid "github.com/ministryofjustice/cloud-platform-grafana-linter/uid"
	u "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"
)

var (
	pull = 0
)

func main() {
	client, ctx := u.Client()

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
