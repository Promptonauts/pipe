package guardrails

import "fmt"

type TokenLimitGuardrail struct {
	MaxTokens int
}

func (g *TokenLimitGuardrail) ID() string {
	return "token-limit"
}

func (g *TokenLimitGuardrail) Phase() Phase {
	return PhasePre
}

func (g *TokenLimitGuardrail) Priority() int {
	return 90
}

func (g *TokenLimitGuardrail) Check(input CheckInput) CheckResult {
	if input.TokenCount > g.MaxTokens {
		return CheckResult{
			Passed : false,
			GuardrailID: g.ID(),
			Message: fmt.Sprintf("token count %d exceeds limit %d", input.TokenCount, g.MaxTokens),
			Action: "block",
		}
	}
	return CheckResult{Passed: true, GuardrailID: g.ID(), Message: "within token limit"}
}