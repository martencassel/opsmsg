# opsmsg

Structured operational messages for Go applications.

## What it does

Turns scattered log messages into a catalog you can version, translate, and maintain separately from your code. Load message templates from YAML, merge catalogs, and dispatch messages through your logging system.

## Usage

```go
import (
    "github.com/martencassel/opsmsg/catalog"
    "github.com/martencassel/opsmsg/dispatcher"
    "github.com/sirupsen/logrus"
)

// Load message catalogs
builtin, _ := catalog.Load("catalog/builtin.yaml")
custom, _ := catalog.Load("catalog/custom.yaml")
merged := catalog.Merge(builtin, custom)

// Setup dispatcher
logger := logrus.New()
d := dispatcher.NewLogrusDispatcher(logger)

// Create and dispatch a message
msg := merged.New("SRV001", map[string]string{"port": "8080"})
d.Dispatch(ctx, msg)
```

## Message catalog format

```yaml
- id: SRV001
  severity: INFO
  text: "Server starting on port {port}"
  help: "Cause: Application startup initiated. Recovery: None required."
  replies: []
```

Each message has an ID, severity level, text template with placeholders, help text explaining cause and recovery, and optional reply suggestions.

## Structure

- `message/` - Message types and severity levels
- `catalog/` - YAML loading and merging
- `dispatcher/` - Output interfaces (logrus, custom formatters)
- `examples/` - Working examples

## Examples

See `examples/` for:
- `todo-app` - Web server with custom and builtin messages
- `message-viewer` - Browse catalog messages
- `formatter-demo` - Custom formatting

Run an example:
```bash
cd examples/todo-app
go run .
```

## Why

Better than string literals scattered in code. Better than i18n when you need operational context. Useful when messages need approval, translation, or runbook links.
