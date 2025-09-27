package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/complytime/compliance-to-policy-plugins/opa-plugin/server"

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
	config *server.Config
}

func NewPlugin() *Plugin {
	return &Plugin{
		config: &server.Config{},
	}
}

func (p *Plugin) Configure(_ context.Context, m map[string]string) error {
	if err := mapstructure.Decode(m, &p.config); err != nil {
		return errors.New("error decoding configuration")
	}
	p.config.Complete()
	return p.config.Validate()
}

func (p *Plugin) Generate(ctx context.Context, pl policy.Policy) error {
	composer := server.NewComposer(p.config.PolicyTemplates, p.config.PolicyOutput)
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
	panic("Implement me")
}
