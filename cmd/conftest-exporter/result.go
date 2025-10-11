package main

import (
	"encoding/json"
	"errors"
	"time"

	ocsf "github.com/Santiago-Labs/go-ocsf/ocsf/v1_5_0"
	"github.com/complytime/complybeacon/proofwatch"
)

// ConftestFinding mirrors the structure of a single "warning" or "failure" entry.
type ConftestFinding struct {
	Message   string                 `json:"message"`
	Policy    string                 `json:"policy"`
	Level     string                 `json:"level"` // "warning" or "failure"
	Namespace string                 `json:"namespace"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ConftestFileResult mirrors the structure of one file's results in the Conftest JSON output.
type ConftestFileResult struct {
	FileName string            `json:"fileName"`
	Warnings []ConftestFinding `json:"warnings"`
	Failures []ConftestFinding `json:"failures"`
}

func (f ConftestFinding) ToOCSF(filename string) (proofwatch.Evidence, error) {
	classUID := 6007
	categoryUID := 6
	categoryName := "Application Activity"
	className := "Scan Activity"
	completedScan := 60070

	// Map operation to OCSF activity type
	var activityID int
	var activityName string
	var typeName string

	vendorName := "conftest"
	productName := "conftest"
	action := "observed"
	actionId := int32(3)
	status, statusID := mapStatus(f)
	numFiles := int32(1)
	unknown := "unknown"
	unknownID := int32(0)

	uid := "conftest-exporter"
	activity := ocsf.ScanActivity{
		ActivityId:   int32(activityID),
		ActivityName: &activityName,
		CategoryName: &categoryName,
		CategoryUid:  int32(categoryUID),
		ClassName:    &className,
		ClassUid:     int32(classUID),
		Status:       &status,
		StatusId:     &statusID,
		Severity:     &unknown,
		SeverityId:   unknownID,
		NumFiles:     &numFiles,
		Message:      &f.Message, // Include the violation message from the finding
		Metadata: ocsf.Metadata{
			Uid: &uid,
			Product: ocsf.Product{
				Name:       &productName,
				VendorName: &vendorName,
			},
			LogProvider: &productName,
		},
		Time:     time.Now().UnixMilli(),
		TypeName: &typeName,
		TypeUid:  int64(completedScan),
	}

	policyData, err := json.Marshal(f)
	if err != nil {
		return proofwatch.Evidence{}, err
	}
	policyDataStr := string(policyData)

	checkId, ok := f.Metadata["short_name"]
	if !ok {
		return proofwatch.Evidence{}, errors.New("expected short_name in metadata")
	}

	checkIdStr, ok := checkId.(string)
	if !ok {
		return proofwatch.Evidence{}, errors.New("expected short_name value to be a string")
	}

	policy := ocsf.Policy{
		Uid:  &checkIdStr,
		Data: &policyDataStr,
	}

	evidenceEvent := proofwatch.Evidence{
		ScanActivity: activity,
		Policy:       policy,
		Action:       &action,
		ActionID:     &actionId,
	}

	return evidenceEvent, nil
}

func mapStatus(f ConftestFinding) (string, int32) {
	return "failure", 2
}
