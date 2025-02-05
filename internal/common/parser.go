package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func LoadConfigFromJSON(configPath string, keys []string) ([]byte, error) {
	var result map[string]any

	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	if err := json.Unmarshal(jsonFile, &result); err != nil {
		return nil, err
	}

	for _, key := range keys {
		if value, ok := result[key]; ok {
			if err := ParseAndSetDuration(key, value, result); err != nil {
				return nil, err
			}
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return jsonBytes, err
}

func ParseAndSetDuration(key string, val any, result map[string]any) error {
	strVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("expected a string value for key '%s'", key)
	}

	duration, err := time.ParseDuration(strVal)
	if err != nil {
		return fmt.Errorf("invalid duration format for '%s': %w", key, err)
	}
	result[key] = duration

	return nil
}
