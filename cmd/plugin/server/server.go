package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/oscal-compass/compliance-to-policy-go/v2/logging"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
)

var (
	_           policy.Provider = (*Plugin)(nil)
	logger      hclog.Logger    = logging.NewPluginLogger()
	regoVersion                 = ast.RegoV1
)

func Logger() hclog.Logger {
	return logger
}

type Plugin struct {
	config     *Config
	lokiClient *LokiClient
}

func NewPlugin() *Plugin {
	return &Plugin{
		config: &Config{},
	}
}

func (p *Plugin) Configure(_ context.Context, m map[string]string) error {
	if err := mapstructure.Decode(m, &p.config); err != nil {
		return errors.New("error decoding configuration")
	}

	// Load configuration from environment variables
	p.config.LoadFromEnv()

	// Initialize Loki client - prioritize Grafana Cloud, fall back to local Loki
	if p.config.GrafanaCloudEndpoint != "" {
		// Use Grafana Cloud endpoint with authentication
		p.lokiClient = NewLokiClientWithAuth(
			p.config.GrafanaCloudEndpoint,
			p.config.GrafanaCloudInstanceID,
			p.config.GrafanaCloudAPIKey,
		)
		logger.Info("Initialized Loki client with Grafana Cloud endpoint")
	} else if p.config.LokiURL != "" {
		p.lokiClient = NewLokiClient(p.config.LokiURL)
		logger.Info("Initialized Loki client with local endpoint")
	}

	return p.config.Validate()
}

func (p *Plugin) Generate(ctx context.Context, pl policy.Policy) error {
	composer := NewComposer(p.config.PolicyTemplates, p.config.PolicyOutput)
	if err := composer.GeneratePolicySet(pl, *p.config); err != nil {
		return fmt.Errorf("error generating policies: %w", err)
	}

	if p.config.Bundle != "" {
		logger.Info(fmt.Sprintf("Creating policy bundle at %s", p.config.Bundle))
		if err := composer.Bundle(context.Background(), *p.config); err != nil {
			return fmt.Errorf("error creating policy bundle: %w", err)
		}
	}
	return nil
}

func (p *Plugin) GetResults(ctx context.Context, pl policy.Policy) (policy.PVPResult, error) {
	logger.Info("GetResults called", "policy_length", len(pl))

	// Initialize result structure
	result := policy.PVPResult{
		ObservationsByCheck: []policy.ObservationByCheck{},
		Links:               []policy.Link{},
	}

	// Check if Loki client is configured
	if p.lokiClient == nil {
		logger.Warn("Loki client not configured, returning empty result")
		return result, nil
	}

	// Process each rule set in the policy
	for _, ruleset := range pl {
		for _, check := range ruleset.Checks {
			policyID := check.ID

			logger.Info("Querying Loki for policy", "policy_id", policyID)

			// Query Loki for logs related to this policy
			logEntries, err := p.lokiClient.QueryLogs(ctx, policyID, 100)
			if err != nil {
				logger.Error("Failed to query Loki", "error", err, "policy_id", policyID)
				continue // Continue with other policies even if one fails
			}

			logger.Info("Retrieved log entries", "count", len(logEntries), "policy_id", policyID)

			// Process log entries into evidence links
			var evidenceLinks []policy.Link
			for _, entry := range logEntries {
				// Create evidence link from log entry
				evidenceLink := policy.Link{
					Description: fmt.Sprintf("Evidence from Loki log entry at %s", entry.Timestamp),
					Href: fmt.Sprintf("%s/loki/api/v1/query_range?query={policy_id=\"%s\"}&start=%s&end=%s",
						p.lokiClient.baseURL, policyID, entry.Timestamp, entry.Timestamp),
				}
				evidenceLinks = append(evidenceLinks, evidenceLink)
			}

			// Create observation for this rule set
			if len(evidenceLinks) > 0 {
				observation := policy.ObservationByCheck{
					Title:             fmt.Sprintf("Policy Check for %s", policyID),
					Description:       fmt.Sprintf("Evidence collected from Loki for policy %s", policyID),
					CheckID:           policyID,
					Methods:           []string{"TEST"},
					Subjects:          []policy.Subject{},
					Collected:         time.Now(),
					RelevantEvidences: evidenceLinks,
					Props: []policy.Property{
						{Name: "source", Value: "loki"},
						{Name: "log_count", Value: fmt.Sprintf("%d", len(logEntries))},
						{Name: "policy_id", Value: policyID},
					},
				}

				// Add subjects for each log entry
				for i, entry := range logEntries {
					subject := policy.Subject{
						Title:       fmt.Sprintf("Log Entry %d", i+1),
						Type:        "resource",
						ResourceID:  fmt.Sprintf("log-%s-%d", policyID, i),
						Result:      mapResults(entry.Labels),
						EvaluatedOn: time.Now(),
						Reason:      "Evidence found in Loki logs",
						Props: []policy.Property{
							{Name: "timestamp", Value: entry.Timestamp},
							{Name: "message", Value: entry.Message},
							{Name: "labels", Value: fmt.Sprintf("%v", entry.Labels)},
						},
					}
					observation.Subjects = append(observation.Subjects, subject)
				}

				result.ObservationsByCheck = append(result.ObservationsByCheck, observation)
			}
		}
	}

	logger.Info("GetResults completed",
		"observations_count", len(result.ObservationsByCheck),
		"total_links", len(result.Links))

	return result, nil
}

func mapResults(labels map[string]string) policy.Result {
	status, ok := labels["compliance_status"]
	if !ok {
		return policy.ResultError
	}
	switch status {
	case "Fail":
		return policy.ResultFail
	case "Pass":
		return policy.ResultPass
	default:
		return policy.ResultError
	}
}
