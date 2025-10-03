package server

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// Required
	PolicyResults string `mapstructure:"policy-results"`

	// Optionally bundle local policy
	Bundle         string `mapstructure:"bundle"`
	BundleRevision string `mapstructure:"bundle-revision"`

	// Optional if building locally
	PolicyTemplates string `mapstructure:"policy-templates"`
	PolicyOutput    string `mapstructure:"policy-output"`

	// Loki Client Config
	LokiURL string `mapstructure:"loki-url"`
}

func (c *Config) Validate() error {
	var errs []error
	if err := checkPath(&c.PolicyResults); err != nil {
		errs = append(errs, err)
	}

	if c.PolicyTemplates != "" {
		if err := checkPath(&c.PolicyOutput); err != nil {
			errs = append(errs, err)
		}

		if err := checkPath(&c.PolicyTemplates); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func checkPath(path *string) error {
	if path != nil && *path != "" {
		cleanedPath := filepath.Clean(*path)
		path = &cleanedPath
		_, err := os.Stat(*path)
		if err != nil {
			return fmt.Errorf("path %q: %w", *path, err)
		}
	}
	return nil
}
