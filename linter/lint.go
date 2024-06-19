package linter

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/google/go-github/v62/github"
	"github.com/grafana/dashboard-linter/lint"
)

var (
	lintConfigFlag  string
	lintVerboseFlag bool
	lintAutofixFlag bool
	lintStrictFlag  bool
)

func ExtractJsonFromYamlFile(file *github.CommitFile) error {
	fileName := file.Filename
	exec.Command("sh", "-c", "yq e '.data[]' "+*fileName+" > dashboard.json").Run()
	// check last command's exit status and if file was created
	if _, err := os.Stat("dashboard.json"); os.IsNotExist(err) {
		return fmt.Errorf("failed to create dashboard.json")
	}
	return nil
}

func LintJsonFile(filename string) (*lint.ResultSet, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	dashboard, err := lint.NewDashboard(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard: %v", err)
	}

	// if no config flag was passed, set a default path of a .lint file in the dashboards directory
	if lintConfigFlag == "" {
		lintConfigFlag = path.Join(path.Dir(filename), ".lint")
	}

	config := lint.NewConfigurationFile()
	if err := config.Load(lintConfigFlag); err != nil {
		return nil, fmt.Errorf("failed to load lint config: %v", err)
	}
	config.Verbose = lintVerboseFlag
	config.Autofix = lintAutofixFlag

	rules := lint.NewRuleSet()
	results, err := rules.Lint([]lint.Dashboard{dashboard})
	if err != nil {
		return nil, fmt.Errorf("failed to lint dashboard: %v", err)
	}

	if config.Autofix {
		changes := results.AutoFix(&dashboard)
		if changes > 0 {
			_, err = fmt.Println(dashboard, filename, buf)
			if err != nil {
				return nil, fmt.Errorf("failed to write dashboard: %v", err)
			}
		}
	}

	results.Configure(config)

	if lintStrictFlag && results.MaximumSeverity() >= lint.Warning {
		return nil, fmt.Errorf("there were linting errors, please see previous output")
	}
	return results, nil
}
