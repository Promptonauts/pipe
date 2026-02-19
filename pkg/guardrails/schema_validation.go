package guardrails

import (
	"encoding/json"
)

type SchemaValidationGuardrail struct {
	ExpectedFields []string
}

func (g *SchemaValidationGuardrail) ID() string {
	return "output-schema-validation"
}

func (g *SchemaValidationGuardrail) Phase() Phase {
	return PhasePost
}

func (g *SchemaValidationGuardrail) Priority() int {
	return 50
}

func (g *SchemaValidationGuardrail) Check(input CheckInput) CheckResult {
	if input.Output == "" {
		return CheckResult{
			Passed:      false,
			GuardrailID: g.ID(),
			Message:     "output is empty",
			Action:      "warn",
		}
	}

	if len(g.ExpectedFields) == 0 {
		return CheckResult{Passed: true, GuardrailID: g.ID(), Message: "no schema to validate"}
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(input.Output), &parsed); err != nil {
		return CheckResult{
			Passed:      false,
			GuardrailID: g.ID(),
			Message:     "output is not valid JSON",
			Action:      "warn",
		}
	}

	for _, field := range g.ExpectedFields {
		if _, ok := parsed[field]; !ok {
			return CheckResult{
				Passed:      false,
				GuardrailID: g.ID(),
				Message:     "missing required field: " + field,
				Action:      "warn",
			}
		}
	}
	return CheckResult{Passed: true, GuardrailID: g.ID(), Message: "schema valid"}
}
