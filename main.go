package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/martian/log"
	l "github.com/ministryofjustice/cloud-platform-grafana-linter/linter"
	u "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"
	githubactions "github.com/sethvargo/go-githubactions"
)

var (
	token          = os.Getenv("AUTH_TOKEN")
	ref            = os.Getenv("GITHUB_REF")
	repo           = os.Getenv("GITHUB_REPOSITORY")
	kubeConfigPath = os.Getenv("KUBE_CONFIG_PATH")
	prf            *u.YAML
	json           *u.JSON
	configMaps     *u.ConfigMaps
	uids           []string
)

func main() {
	flag.Parse()

	// chech env vars are set
	if token == "" {
		log.Errorf("AUTH_TOKEN is not set")
	}

	if ref == "" {
		log.Errorf("GITHUB_REF is not set")
	}

	if repo == "" {
		log.Errorf("GITHUB_REPOSITORY is not set")
	}

	if kubeConfigPath == "" {
		log.Errorf("KUBE_CONFIG_PATH is not set")
	}

	// gather config maps from all namespaces with label grafana_dashboard
	kube_client := u.ClientSet(kubeConfigPath)

	configList, err := u.NamespaceConfigMaps(kube_client)
	if err != nil {
		log.Errorf("Error fetching config maps: %v\n", err)
	}

	for _, configMap := range configList {
		for _, items := range configMap.Items {
			for _, value := range items.Data {
				cm, err := u.ExtractJson(value, "uid")
				if err != nil {
					fmt.Printf("Error extracting json: %v\n", err)
				}
				configMaps, _ = cm.(*u.ConfigMaps)

				// create an array of configMaps.UID
				uids = []string{configMaps.UID}
			}
		}
	}

	// get pull request files
	githubrefS := strings.Split(ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	client := u.GitHubClient(token)

	files, _, err := u.GetPullRequestFiles(client, owner, repoName, pull)
	if err != nil {
		log.Errorf("Error fetching pull request files: %v\n", err)
	}

	for _, file := range files {
		if u.SelectFile(pull, file) != nil {
			decodeContent, err := u.GetFileContent(client, file, owner, repoName)
			if err != nil {
				log.Errorf("Error fetching file content: %v\n", err)
			}

			prf, err = u.ExtractYaml(decodeContent)
			if err != nil {
				log.Errorf("Error extracting yaml: %v\n", err)
			}

			for _, v := range prf.Data {
				j, err := u.ExtractJson(v, "linter")
				if err != nil {
					log.Errorf("Error extracting json: %v\n", err)
				}
				json, _ = j.(*u.JSON)
				fmt.Println("UID: ", json.UID)

				for _, uid := range uids {
					if uid == json.UID {
						log.Infof("UID exists in cluster")
						// github output for actions
						githubactions.New().SetOutput("uid_exists", "true")
						githubactions.New().SetOutput("uid_message", "UID: "+json.UID+" already exists in cluster")
						os.Exit(0)
					} else {
						fmt.Println("UID does not exist in cluster")
					}
				}
			}

			fmt.Printf("\nTitle: %s\n", prf.Metadata.Name)
			fmt.Printf("Namespace: %s\n", prf.Metadata.Namespace)

			// print data interface
			for k, v := range prf.Data {
				results, err := l.LintJsonFile(k, []byte(v.(string)))
				if err != nil {
					log.Errorf("Error linting json file: %v\n", err)
				}

				results.ReportByRule()

				githubactions.New().SetOutput("uid_exists", "false")
				githubactions.New().SetOutput("uid_message", "UID: "+json.UID+" does not exist in cluster")

			}

		} else {
			continue
		}
	}

}
