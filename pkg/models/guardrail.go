package models

type GuardrailSpec struct {
	Description string                 `yaml:"description" json:"description"`
	Type        string                 `yaml:"type" json:"type"`
	Phase       string                 `yaml:"phase" json:"phase"`
	Config      map[string]interface{} `yaml:"config" json:"config"`
	Action      string                 `yaml:"action" json:"action"`
	Priority    int                    `yaml:"priority" json:"priority"`
}