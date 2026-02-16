package observability

import (
	"time"

	"github.com/google/uuid"
)

type Span struct {
	TraceID   string            `json:"traceID"`
	SpanID    string            `json:"spanID"`
	ParentID  string            `json:"parentID,omitempty"`
	Operation string            `json:"operation"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime,omitempty"`
	Duration  time.Duration     `json:"duration,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Status    string            `json:"status"`
	Events    []SpanEvent       `json:"events,omitempty"`
}

type SpanEvent struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type Tracer struct {
	logger *Logger
}

func NewTracer(logger *Logger) *Tracer {
	return &Tracer{logger: logger.With("tracer")}
}

func (t *Tracer) StartSpan(operation string, traceID string) *Span {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	span := &Span{
		TraceID:   traceID,
		SpanID:    uuid.New().String(),
		Operation: operation,
		StartTime: time.Now().UTC(),
		Status:    "started",
		Tags:      make(map[string]string),
	}
	t.logger.Debug("span started", "traceID", traceID, "spanID", span.SpanID, "op", operation)
	return span
}

func (t *Tracer) EndSpan(span *Span, status string) {
	span.EndTime = time.Now().UTC()
	span.Duration = span.EndTime.Sub(span.StartTime)
	span.Status = status
	t.logger.Debug("span ended", "traceId", span.TraceID, "spanId", span.SpanID, "op", span.Operation, "durationMs", span.Duration.Milliseconds(), "status", status)
}
