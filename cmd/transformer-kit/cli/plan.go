package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

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

			if strings.TrimSpace(targetComponent) == "" {
				return fmt.Errorf("target-component is required: component is missing a title")
			}

			builder = builder.AddTargetComponent(targetComponent, componentType, layer2Catalog, nil)
			builder = builder.AddValidationComponent(layer4Plan)

			cleanedPolicyPath := filepath.Clean(policyPath)
			var layer3Policy layer3.PolicyDocument
			if err := layer3Policy.LoadFile(fmt.Sprintf("file://%s", cleanedPolicyPath)); err != nil {
				return err
			}
			compDef := builder.Build()
			compDef = removeEmptyProps(compDef)

			var found bool
			for _, guidance := range layer3Policy.GuidanceReferences {
				if guidanceRef == guidance.ReferenceId {
					ap, err := transformers.ComponentDefinitionsToAssessmentPlan(cmd.Context(), []oscalTypes.ComponentDefinition{compDef}, guidance.ReferenceId)
					if err != nil {
						return err
					}
					ap = removeEmptyPropsFromAssessmentPlan(ap)
					ap = fillEmptyComponentTitles(ap)
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
	flags.StringVarP(&catalogPath, "catalog-path", "c", "./governance/catalogs/cnscc.yaml", "Path to L2 Catalog to transform")
	flags.StringVarP(&evaluationsPath, "evaluation-path", "e", "./governance/plans/cnscc.yaml", "Path to Layer 4 Evaluation Plan to transform")
	flags.StringVarP(&targetComponent, "target-component", "t", "", "Title for target component for evaluation")
	flags.StringVar(&componentType, "component-type", "software", "Component type (based on valid OSCAL component types)")
	flags.StringVarP(&policyPath, "policy-path", "p", "./governance/policy.yaml", "Path to Layer 3 policy")
	flags.StringVarP(&guidanceRef, "guidance-reference", "r", "", "Guidance reference to tailor the plan to")
	return command
}

// removeEmptyProps removes props with empty values from the component definition
// to comply with OSCAL schema validation that requires prop values to match
// the pattern '^\\S(.*\\S)?$' (non-empty, no leading/trailing whitespace)
func removeEmptyProps(compDef oscalTypes.ComponentDefinition) oscalTypes.ComponentDefinition {
	if compDef.Components == nil {
		return compDef
	}

	components := *compDef.Components
	for i := range components {
		if components[i].Props == nil {
			continue
		}

		props := *components[i].Props
		if len(props) == 0 {
			continue
		}

		// Use reflection to get the element type and create a new slice
		propType := reflect.TypeOf(props[0])
		filteredSlice := reflect.MakeSlice(reflect.SliceOf(propType), 0, len(props))

		for _, prop := range props {
			// Use reflection to access the Value field
			propValue := reflect.ValueOf(prop)
			valueField := propValue.FieldByName("Value")
			if !valueField.IsValid() || valueField.Kind() != reflect.String {
				// If we can't access Value, keep the prop
				filteredSlice = reflect.Append(filteredSlice, propValue)
				continue
			}

			strValue := valueField.String()
			trimmedValue := strings.TrimSpace(strValue)
			if trimmedValue != "" {
				// Create a new prop with trimmed value
				newProp := reflect.New(propType).Elem()
				newProp.Set(propValue)
				newProp.FieldByName("Value").SetString(trimmedValue)
				filteredSlice = reflect.Append(filteredSlice, newProp)
			}
		}

		// Props is already a pointer to slice, so we need to create a pointer to the filtered slice
		filteredSliceValue := reflect.New(reflect.SliceOf(propType))
		filteredSliceValue.Elem().Set(filteredSlice)
		// Use unsafe pointer conversion or direct assignment via reflection
		// Since we know the type matches, we can use a type assertion
		components[i].Props = filteredSliceValue.Interface().(*[]oscalTypes.Property)
	}
	compDef.Components = &components

	return compDef
}

// removeEmptyPropsFromAssessmentPlan removes props with empty values from assessment plan components
// to comply with OSCAL schema validation that requires prop values to match
// the pattern '^\\S(.*\\S)?$' (non-empty, no leading/trailing whitespace)
func removeEmptyPropsFromAssessmentPlan(ap *oscalTypes.AssessmentPlan) *oscalTypes.AssessmentPlan {
	if ap == nil || ap.AssessmentAssets == nil {
		return ap
	}

	if ap.AssessmentAssets.Components == nil {
		return ap
	}

	components := *ap.AssessmentAssets.Components
	for i := range components {
		if components[i].Props == nil {
			continue
		}

		props := *components[i].Props
		if len(props) == 0 {
			continue
		}

		// Use reflection to get the element type and create a new slice
		propType := reflect.TypeOf(props[0])
		filteredSlice := reflect.MakeSlice(reflect.SliceOf(propType), 0, len(props))

		for _, prop := range props {
			// Use reflection to access the Value field
			propValue := reflect.ValueOf(prop)
			valueField := propValue.FieldByName("Value")
			if !valueField.IsValid() || valueField.Kind() != reflect.String {
				// If we can't access Value, keep the prop
				filteredSlice = reflect.Append(filteredSlice, propValue)
				continue
			}

			strValue := valueField.String()
			trimmedValue := strings.TrimSpace(strValue)
			if trimmedValue != "" {
				// Create a new prop with trimmed value
				newProp := reflect.New(propType).Elem()
				newProp.Set(propValue)
				newProp.FieldByName("Value").SetString(trimmedValue)
				filteredSlice = reflect.Append(filteredSlice, newProp)
			}
		}

		// Props is already a pointer to slice, so we need to create a pointer to the filtered slice
		filteredSliceValue := reflect.New(reflect.SliceOf(propType))
		filteredSliceValue.Elem().Set(filteredSlice)
		// Use unsafe pointer conversion or direct assignment via reflection
		// Since we know the type matches, we can use a type assertion
		components[i].Props = filteredSliceValue.Interface().(*[]oscalTypes.Property)
	}
	ap.AssessmentAssets.Components = &components

	return ap
}

// fillEmptyComponentTitles fills any empty component titles in the assessment plan with "REPLACE_ME"
func fillEmptyComponentTitles(ap *oscalTypes.AssessmentPlan) *oscalTypes.AssessmentPlan {
	if ap == nil {
		return ap
	}

	// Fill empty titles in assessment-assets.components
	if ap.AssessmentAssets != nil && ap.AssessmentAssets.Components != nil {
		components := *ap.AssessmentAssets.Components
		for i := range components {
			if strings.TrimSpace(components[i].Title) == "" {
				title := "REPLACE_ME"
				components[i].Title = title
			}
		}
		ap.AssessmentAssets.Components = &components
	}

	if ap.LocalDefinitions != nil && ap.LocalDefinitions.Components != nil {
		components := *ap.LocalDefinitions.Components
		for i := range components {
			if strings.TrimSpace(components[i].Title) == "" {
				title := "REPLACE_ME"
				components[i].Title = title
			}
		}
		ap.LocalDefinitions.Components = &components
	}

	// Fill empty titles in assessment-assets.assessment-platforms
	if ap.AssessmentAssets != nil && ap.AssessmentAssets.AssessmentPlatforms != nil {
		platforms := ap.AssessmentAssets.AssessmentPlatforms
		for i := range platforms {
			if strings.TrimSpace(platforms[i].Title) == "" {
				title := "REPLACE_ME"
				platforms[i].Title = title
			}
		}
		ap.AssessmentAssets.AssessmentPlatforms = platforms
	}

	return ap
}
