package uid

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v62/github"
)

func SearchCode(pull int, file *github.CommitFile) (string, error) {
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
	return "", fmt.Errorf("error: UID not found in PRs")
}
