package main

import (
	"fmt"
	"os"

	"github.com/ministryofjustice/cloud-platform-grafana-linter/uid/cluster"
	"github.com/ministryofjustice/cloud-platform-grafana-linter/uid/github"
)

func main() {
	client, ctx := github.Client()
	uid, err := github.SearchCode(client, ctx)
	if err != nil {
		fmt.Printf("Error searching code: %v\n", err)
		os.Exit(1)
	}

	clientset := cluster.ClientSet()
	configMaps, err := cluster.SearchNamespacesForConfigMaps(clientset)
	if err != nil {
		fmt.Printf("Error getting configmaps: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %v configmaps\n", len(configMaps))
	for _, configMap := range configMaps {
		d := cluster.SearchConfigMapsData(clientset, configMap, uid)
		if d != "" {
			fmt.Printf("%v", d)
		}
	}
}
