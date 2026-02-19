package guardrails

import (
	"strings"
)

type PromptInjectionGuardrail struct{}

func (g *PromptInjectionGuardrail) ID() string {
	return "prompt-injection"
}

func (g *PromptInjectionGuardrail) Phase() Phase {
	return PhasePre
}

func (g *PromptInjectionGuardrail) Priority() int {
	return 100
}

var injectionPatterns = []string{
	"ignore previous instructions",
	"ignore all previous",
	"disregard above",
	"forget your instructions",
	"you are now",
	"act as if",
	"pretend you are",
	"override your",
	"new instructions:",
	"system prompt:",
	"ignore the above",
	"do not follow",
	"bypass your",
	"reveal your system",
	"show me your prompt",
	"what is your system prompt",
	//...........so on . can be added specifically  , depends on the use case.
}

func (g *PromptInjectionGuardrail) Check(input CheckInput) CheckResult {
	lower := strings.ToLower(input.Prompt)

	for _, pattern := range injectionPatterns {
		if strings.Contains(lower, pattern) {
			return CheckResult{
				Passed:      false,
				GuardrailID: g.ID(),
				Message:     "potential prompt injection detected: matched pattern '" + pattern + "'",
				Action:      "block",
			}
		}
	}
	return CheckResult{
		Passed:      true,
		GuardrailID: g.ID(),
		Message:     "no injection detected",
		Action:      "",
	}
}
