package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func loadFiles() (fileConfig, error) {
	env := strings.ToLower(stringEnv("OUTBOUND_ENV", "development"))
	configDir := stringEnv("OUTBOUND_CONFIG_DIR", "configs")

	merged := map[string]any{}
	for _, path := range []string{
		filepath.Join(configDir, "default.yaml"),
		filepath.Join(configDir, env+".yaml"),
	} {
		if err := mergeYAMLFile(merged, path); err != nil {
			return fileConfig{}, err
		}
	}

	raw, err := yaml.Marshal(merged)
	if err != nil {
		return fileConfig{}, fmt.Errorf("marshal merged config: %w", err)
	}

	var cfg fileConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return fileConfig{}, fmt.Errorf("unmarshal merged config: %w", err)
	}
	return cfg, nil
}

func mergeYAMLFile(target map[string]any, path string) error {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	var source map[string]any
	if err := yaml.Unmarshal(content, &source); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	mergeMap(target, source)
	return nil
}

func mergeMap(target, source map[string]any) {
	for key, sourceValue := range source {
		sourceMap, sourceIsMap := sourceValue.(map[string]any)
		targetMap, targetIsMap := target[key].(map[string]any)
		if sourceIsMap && targetIsMap {
			mergeMap(targetMap, sourceMap)
			continue
		}
		target[key] = sourceValue
	}
}
