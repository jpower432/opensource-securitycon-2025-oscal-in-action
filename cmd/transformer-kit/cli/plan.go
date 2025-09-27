package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/complytime/gemara2oscal/component"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
	"github.com/ossf/gemara/layer2"
	"github.com/ossf/gemara/layer3"
	"github.com/ossf/gemara/layer4"
	"github.com/spf13/cobra"
)

func NewPlanCommand() *cobra.Command {
	var catalogPath, targetComponent, componentType, evaluationsPath, policyPath, guidanceRef string

	command := &cobra.Command{
		Use:   "plan",
		Short: "Transform Gemara governance artifacts to an OSCAL Assessment Plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := component.NewDefinitionBuilder("GitHub Repository", "0.1.0")

			cleanedCatalogPath := filepath.Clean(catalogPath)
			catalogData, err := os.ReadFile(cleanedCatalogPath)
			if err != nil {
				return err
			}

			var layer2Catalog layer2.Catalog
			if err := layer2Catalog.LoadFile(fmt.Sprintf("file://%s", cleanedCatalogPath)); err != nil {
				return err
			}
			err = yaml.Unmarshal(catalogData, &layer2Catalog)
			if err != nil {
				return err
			}

			cleanedPlanPath := filepath.Clean(evaluationsPath)
			planBytes, err := os.ReadFile(cleanedPlanPath)
			if err != nil {
				return err
			}
			var layer4Plan layer4.EvaluationPlan
			err = yaml.Unmarshal(planBytes, &layer4Plan)
			if err != nil {
				return err
			}

			builder = builder.AddTargetComponent(targetComponent, componentType, layer2Catalog)
			builder = builder.AddValidationComponent(layer4Plan)

			cleanedPolicyPath := filepath.Clean(policyPath)
			var layer3Policy layer3.PolicyDocument
			if err := layer3Policy.LoadFile(fmt.Sprintf("file://%s", cleanedPolicyPath)); err != nil {
				return err
			}

			for _, ref := range layer3Policy.ControlReferences {
				builder = builder.AddParameterModifiers(ref.ReferenceId, ref.ParameterModifications)
			}
			compDef := builder.Build()

			var found bool
			for _, guidance := range layer3Policy.GuidanceReferences {
				if guidanceRef == guidance.ReferenceId {
					ap, err := transformers.ComponentDefinitionsToAssessmentPlan(cmd.Context(), []oscalTypes.ComponentDefinition{compDef}, guidance.ReferenceId)
					if err != nil {
						return err
					}
					oscalModels := oscalTypes.OscalModels{
						AssessmentPlan: ap,
					}
					compDefData, err := json.MarshalIndent(oscalModels, "", " ")
					if err != nil {
						return err
					}
					_, _ = fmt.Fprintln(os.Stdout, string(compDefData))
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("guidance reference %q does not exist in policy", guidanceRef)
			}

			return nil
		},
	}

	flags := command.Flags()
	flags.StringVarP(&catalogPath, "catalog-path", "c", "./governance/catalogs/osps.yaml", "Path to L2 Catalog to transform")
	flags.StringVarP(&evaluationsPath, "evaluation-path", "e", "./governance/plans/osps.yaml", "Path to Layer 4 Evaluation Plan to transform")
	flags.StringVarP(&targetComponent, "target-component", "t", "", "Title for target component for evaluation")
	flags.StringVar(&componentType, "component-type", "software", "Component type (based on valid OSCAL component types)")
	flags.StringVarP(&policyPath, "policy-path", "p", "./governance/policy.yaml", "Path to Layer 3 policy")
	flags.StringVarP(&guidanceRef, "guidance-reference", "r", "", "Guidance reference to tailor the plan to")
	return command
}
