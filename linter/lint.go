package linter

import (
	"fmt"

	"github.com/grafana/dashboard-linter/lint"
)

var (
	lintConfigFlag  string
	lintVerboseFlag bool
	lintAutofixFlag bool
	lintStrictFlag  bool
)

func LintJsonFile(key string, value []byte) (*lint.ResultSet, error) {
	dashboard, err := lint.NewDashboard(value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dashboard: %v", err)
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
			_, err = fmt.Println(dashboard, key, value)
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
