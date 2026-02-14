package observability

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	component string
}

func NewLogger(component string) *Logger {
	return &Logger{component: component}
}

type logEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

func (l *Logger) log(level, msg string, fields map[string]interface{}) {
	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Component: l.component,
		Message:   msg,
		Fields:    fields,
	}
	data, _ := json.Marshal(entry)
	fmt.Fprintln(os.Stderr, string(data))
}

func (l *Logger) Info(msg string, kv ...interface{}) {
	l.log("INFO", msg, toMap(kv))
}

func (l *Logger) Error(msg string, kv ...interface{}) {
	l.log("ERROR", msg, toMap(kv))
}

func (l *Logger) Warn(msg string, kv ...interface{}) {
	l.log("WARN", msg, toMap(kv))
}

func (l *Logger) Debug(msg string, kv ...interface{}) {
	l.log("DEBUG", msg, toMap(kv))
}

func (l *Logger) With(component string) *Logger {
	return &Logger{component: l.component + "." + component}
}

func toMap(kv []interface{}) map[string]interface{} {
	if len(kv) == 0 {
		return nil
	}
	m := make(map[string]interface{})
	for i := 0; i < len(kv)-1; i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", kv[i])
		}
		m[key] = kv[i+1]
	}
	return m
}
