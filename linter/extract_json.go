package linter

import (
	"fmt"
	"os"

	"github.com/google/go-github/v57/github"
)

func ExtractJsonFromPullRequestFile(files []*github.CommitFile) (string, error) {
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
	return "dashboard.json: created", nil
}
