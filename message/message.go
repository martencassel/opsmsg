package message

import "time"

type Severity string

const (
    Info     Severity = "INFO"
    Warn     Severity = "WARN"
    Error    Severity = "ERROR"
    Critical Severity = "CRITICAL"
)

type Message struct {
    ID        string
    Severity  Severity
    Text      string
    Context   map[string]string
    Timestamp time.Time
    Help      string
    Replies   []string
}
