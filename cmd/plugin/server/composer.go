package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/open-policy-agent/opa/v1/compile"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	cp "github.com/otiai10/copy"
)

// Adapted from: https://github.com/complytime/compliance-to-policy-plugins/blob/main/opa-plugin/server/composer.go
// SPDX-License-Identifier: Apache-2.0

type Composer struct {
	policiesTemplates string
	policyOutput      string
}

func NewComposer(policiesTemplates string, output string) *Composer {
	return &Composer{
		policiesTemplates: policiesTemplates,
		policyOutput:      output,
	}
}

func (c *Composer) GetPoliciesDir() string {
	return c.policiesTemplates
}

func (c *Composer) Bundle(ctx context.Context, config Config) error {
	buf := bytes.NewBuffer(nil)

	compiler := compile.New().
		WithRevision(config.BundleRevision).
		WithOutput(buf).
		WithPaths(config.PolicyOutput)

	compiler = compiler.WithRegoVersion(regoVersion)

	err := compiler.Build(ctx)
	if err != nil {
		return err
	}

	out, err := os.Create(config.Bundle)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Composer) GeneratePolicySet(pl policy.Policy, config Config) error {
	var local bool
	outputDir := c.policyOutput
	if config.PolicyTemplates != "" {
		local = true
		outputDir = filepath.Join(c.policyOutput, "policy")
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
	}

	for _, rule := range pl {
		for _, check := range rule.Checks {
			if local {
				// Copy over the check directory
				origfilePath := filepath.Join(c.policiesTemplates, check.ID)
				destfilePath := filepath.Join(outputDir, check.ID)
				if err := cp.Copy(origfilePath, destfilePath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
