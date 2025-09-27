package main

import "github.com/complytime/complybeacon/proofwatch"

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

// TODO: Implement

func (f ConftestFinding) ToOCSF() proofwatch.Evidence {
	return proofwatch.Evidence{}
}
