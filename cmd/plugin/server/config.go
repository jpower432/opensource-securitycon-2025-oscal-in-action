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

	// Grafana Cloud Config
	GrafanaCloudEndpoint   string `mapstructure:"grafana-cloud-endpoint"`
	GrafanaCloudInstanceID string `mapstructure:"grafana-cloud-instance-id"`
	GrafanaCloudAPIKey     string `mapstructure:"grafana-cloud-api-key"`
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	// Load Grafana Cloud credentials from environment variables
	if c.GrafanaCloudEndpoint == "" {
		c.GrafanaCloudEndpoint = os.Getenv("GRAFANA_CLOUD_LOKI_ENDPOINT")
	}
	if c.GrafanaCloudInstanceID == "" {
		c.GrafanaCloudInstanceID = os.Getenv("GRAFANA_CLOUD_INSTANCE_ID")
	}
	if c.GrafanaCloudAPIKey == "" {
		c.GrafanaCloudAPIKey = os.Getenv("GRAFANA_CLOUD_API_KEY")
	}
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

	// Validate Loki configuration
	if c.LokiURL == "" && c.GrafanaCloudEndpoint == "" {
		errs = append(errs, errors.New("either loki-url or grafana-cloud-endpoint must be provided"))
	}

	// Validate Grafana Cloud authentication
	if c.GrafanaCloudEndpoint != "" {
		if c.GrafanaCloudInstanceID == "" || c.GrafanaCloudAPIKey == "" {
			errs = append(errs, errors.New("both grafana-cloud-instance-id and grafana-cloud-api-key must be provided when using grafana-cloud-endpoint"))
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
