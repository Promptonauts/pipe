package models

type ToolSpec struct {
	Description string            `yaml:"description" json:"description"`
	Type        string            `yaml:"type" json:"type"`
	Endpoint    string            `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Command     string            `yaml:"command,omitempty" json:"command,omitempty"`
	Schema      ToolSchema        `yaml:"schema" json:"schema"`
	Config      map[string]string `yaml:"config,omitempty" json:"config,omitempty"`
	Timeout     string            `yaml:"timeout" json:"timeout"`
}

type ToolSchema struct {
	Input  map[string]interface{} `yaml:"input" json:"input"`
	Output map[string]interface{} `yaml:"output" json:"output"`
}