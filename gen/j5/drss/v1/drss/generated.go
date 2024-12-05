package drss

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/j5/drss/v1/drss

import ()

// StepMeta Proto: StepMeta
type StepMeta struct {
	StepId    string   `json:"stepId,omitempty"`
	Name      string   `json:"name,omitempty"`
	Status    string   `json:"status,omitempty"`
	Error     *string  `json:"error,omitempty"`
	DependsOn []string `json:"dependsOn,omitempty"`
}

// StepResult Proto: StepResult
type StepResult struct {
	StepId string  `json:"stepId,omitempty"`
	Status string  `json:"status,omitempty"`
	Error  *string `json:"error,omitempty"`
}

// StepStatus Proto Enum: j5.drss.v1.StepStatus
type StepStatus string

const (
	StepStatus_UNSPECIFIED StepStatus = "UNSPECIFIED"
	StepStatus_BLOCKED     StepStatus = "BLOCKED"
	StepStatus_READY       StepStatus = "READY"
	StepStatus_ACTIVE      StepStatus = "ACTIVE"
	StepStatus_DONE        StepStatus = "DONE"
	StepStatus_FAILED      StepStatus = "FAILED"
)
