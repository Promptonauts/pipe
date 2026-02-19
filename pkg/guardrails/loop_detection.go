package guardrails

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type LoopDetectionGuardrail struct {
	MaxRepeats int
	mu         sync.Mutex
	seen       map[string]int //hash -> count
}

func (g *LoopDetectionGuardrail) ID() string {
	return "loop-detection"
}

func (g *LoopDetectionGuardrail) Phase() Phase {
	return "both"
}

func (g *LoopDetectionGuardrail) Priority() int {
	return 80
}

func (g *LoopDetectionGuardrail) Check(input CheckInput) CheckResult {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.seen == nil {
		g.seen = make(map[string]int)
	}

	content := input.Prompt + "|" + input.Output
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	key := input.ExecutionID + ":" + hash

	g.seen[key]++
	count := g.seen[key]

	if count > g.MaxRepeats {
		return CheckResult{
			Passed:      false,
			GuardrailID: g.ID(),
			Message:     fmt.Sprintf("loop detected: same content repeated %d times (max %d)", count, g.MaxRepeats),
			Action:      "block",
		}
	}
	return CheckResult{Passed: true, GuardrailID: g.ID(), Message: "no loop detected"}
}
