package models

type AgentSpec struct {
	Description string            `yaml:"description" json:"description"`
	Runtime     string            `yaml:"runtime" json:"runtime"`
	Model       ModelConfig       `yaml:"model" json:"model"`
	Tools       []string          `yaml:"tools,omitempty" json:"tools,omitempty"`
	Guardrails  []string          `yaml:"guardrails,omitempty" json:"guardrails,omitempty"`
	Config      map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
	MaxRetries  int               `yaml:"maxRetries" json:"maxRetries"`
}

type ModelConfig struct {
	Provider    string  `yaml:"provider" json:"provider"`
	Name        string  `yaml:"name" json:"name"`
	Temperature float64 `yaml:"temperature" json:"temperature"`
	MaxTokens   int     `yaml:"maxTokens" json:"maxTokens"`
}
