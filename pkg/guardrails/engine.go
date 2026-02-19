package guardrails

import (
	"fmt"

	"github.com/Promptonauts/pipe/pkg/observability"
)

type Phase string

const (
	PhasePre  Phase = "pre"
	PhasePost Phase = "post"
)

type CheckInput struct {
	Prompt      string
	Output      string
	TokenCount  int
	AgentName   string
	ExecutionID string
	StepIndex   int
	Metadata    map[string]interface{}
}

type CheckResult struct {
	Passed      bool
	GuardrailID string
	Message     string
	Action      string // block , warn or log
}

type Guardrail interface {
	ID() string
	Phase() Phase
	Check(input CheckInput) CheckResult
	Priority() int
}

type Engine struct {
	guardrails []Guardrail
	metrics    *observability.MetricsRegistry
	logger     *observability.Logger
}

func NewEngine(metrics *observability.MetricsRegistry, logger *observability.Logger) *Engine {
	e := &Engine{
		metrics: metrics,
		logger:  logger.With("guardrails"),
	}
	e.Register(&PromptInjectionGuardrail{})

	e.Register(&TokenLimitGuardrail{MaxTokens: 4096})

	e.Register(&LoopDetectionGuardrail{MaxRepeats: 3})

	e.Register(NewRateLimiter(100)) //100 requests per minitue
	return e
}

func (e *Engine) Register(g Guardrail) {
	e.guardrails = append(e.guardrails, g)
	e.logger.Info("guardrail registered", "id", g.ID(), "phase", g.Phase())
}

func (e *Engine) RunPre(input CheckInput) ([]CheckResult, error) {
	return e.run(PhasePre, input)
}

func (e *Engine) RunPost(input CheckInput) ([]CheckResult, error) {
	return e.run(PhasePost, input)
}

func (e *Engine) run(phase Phase, input CheckInput) ([]CheckResult, error) {
	var results []CheckResult

	for _, g := range e.guardrails {
		if g.Phase() != phase && g.Phase() != "both" {
			continue
		}
		result := g.Check(input)
		results = append(results, result)

		e.metrics.Counter("guardrail.checks.total").Inc()

		if !result.Passed {

			e.metrics.Counter("guardrail.violations.total").Inc()

			e.metrics.Counter(fmt.Sprintf("gyardrail.violations.%s", g.ID())).Inc()
			e.logger.Warn("guardrail violation",
				"guardrail", g.ID(),
				"action", result.Action,
				"message", result.Message,
				"agent", input.AgentName,
				"execution", input.ExecutionID,
			)
			if result.Action == "block" {
				return results,
					fmt.Errorf("blocked by guardrail %s: %s", g.ID(), result.Message)
			}
		}
	}
	return results, nil
}
