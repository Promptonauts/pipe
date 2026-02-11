package models

import (
	"fmt"
	"time"
)

type ResourceKind string

const (
	KindAgent     ResourceKind = "Agent"
	KindTool      ResourceKind = "Tool"
	KindGuardrail ResourceKind = "Guardrail"
	KindPipeline  ResourceKind = "Pipeline"
	KindExecution ResourceKind = "Execution"
)

func ParseResourceKind(s string) (ResourceKind, error) {
	switch s {
	case "Agent":
		return KindAgent, nil
	case "Tool":
		return KindTool, nil
	case "Guardrail":
		return KindGuardrail, nil
	case "Pipeline":
		return KindPipeline, nil
	case "Execution":
		return KindExecution, nil
	default:
		return "", fmt.Errorf("unknown resource kind: %s", s)
	}
}

type Metadata struct {
	Name      string            `yaml:"name" json:"name"`
	Namespace string            `yaml:"namespace" json:"namespace"`
	Labels    map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Version   string            `yaml:"version" json:"version"`
	UID       string            `yaml:"uid,omitempty" json:"uid,omitempty"`
	CreatedAt time.Time         `yaml:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time         `yaml:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

type ResourceStatus struct {
	State       string            `yaml:"state" json:"state"`
	Health      string            `yaml:"health" json:"health"`
	LastUpdated time.Time         `yaml:"lastUpdated" json:"lastUpdated"`
	Metrics     map[string]int64  `yaml:"metrics,omitempty" json:"metrics,omitempty"`
	Conditions  map[string]string `yaml:"conditions,omitempty" json:"conditions,omitempty"`
}

type GenericResource struct {
	APIVersion string                 `yaml:"apiVersion" json:"apiVersion"`
	Kind       ResourceKind           `yaml:"kind" json:"kind"`
	Metadata   Metadata               `yaml:"metadata" json:"metadata"`
	Spec       map[string]interface{} `yaml:"spec" json:"spec"`
	Status     ResourceStatus         `yaml:"status,omitempty" json:"status,omitempty"`
}

func (r *GenericResource) Key() string {
	return fmt.Sprintf("%s/%s/%s", r.Kind, r.Metadata.Namespace, r.Metadata.Name)
}