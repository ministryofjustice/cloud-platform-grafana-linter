package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

var (
	j  *JSON
	cm *ConfigMaps
)

func ExtractYaml(decodeContent string) (*YAML, error) {
	var configMap YAML

	err := yaml.Unmarshal([]byte(decodeContent), &configMap)
	if err != nil {
		return nil, err
	}

	return &configMap, nil
}

func ExtractJson(value interface{}, s string) (interface{}, error) {
	if s == "uid" {
		err := json.Unmarshal([]byte(value.(string)), &cm)
		if err != nil {
			return nil, err
		}

		return cm, nil
	}

	if s == "linter" {
		err := json.Unmarshal([]byte(value.(string)), &j)
		if err != nil {
			return nil, err
		}

		return j, nil
	}

	return nil, fmt.Errorf("unknown type: %s", s)
}
