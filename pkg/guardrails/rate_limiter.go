package guardrails

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiterGuardrail struct {
	maxPerMinute int
	mu           sync.Mutex
	windows      map[string][]time.Time
}

func NewRateLimiter(maxPerMinute int) *RateLimiterGuardrail {
	return &RateLimiterGuardrail{
		maxPerMinute: maxPerMinute,
		windows:      make(map[string][]time.Time),
	}
}

func (g *RateLimiterGuardrail) ID() string    { return "rate-limiter" }
func (g *RateLimiterGuardrail) Phase() Phase  { return PhasePre }
func (g *RateLimiterGuardrail) Priority() int { return 95 }

func (g *RateLimiterGuardrail) Check(input CheckInput) CheckResult {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-1 * time.Minute)
	key := input.AgentName

	//Pruning the old entries

	var valid []time.Time
	for _, t := range g.windows[key] {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	g.windows[key] = valid

	if len(valid) >= g.maxPerMinute {
		return CheckResult{
			Passed:      false,
			GuardrailID: g.ID(),
			Message:     fmt.Sprintf("rate limit exceeded: %d requests in last minute (max %d)", len(valid), g.maxPerMinute),
			Action:      "block",
		}
	}
	g.windows[key] = append(g.windows[key], now)
	return CheckResult{Passed: true, GuardrailID: g.ID(), Message: "within rate limit"}
}
