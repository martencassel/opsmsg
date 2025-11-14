# Formatter Demo

Demonstrates the built-in IBM-style formatters for Logrus that can be used with the opsmsg dispatcher.

## Features

This example shows how to use:

1. **IBMFormatter** - Full IBM-style format with box borders and colors
2. **SimpleIBMFormatter** - Simplified IBM format without borders
3. **Direct Logrus integration** - Use standard Logrus logger idioms with IBM formatting
4. **Color control** - Enable/disable ANSI colors

## Running

```bash
go run main.go
```

## Usage Patterns

### 1. With Dispatcher (Recommended)

```go
logger := logrus.New()
logger.SetFormatter(&dispatcher.IBMFormatter{
    Width: 80,
})

d := dispatcher.NewLogrusDispatcher(logger)

// Use catalog messages
msg := catalog.New("SRV002", map[string]string{
    "port": "8080",
    "error": "address already in use",
})
d.Dispatch(context.Background(), msg)
```

### 2. Direct Logrus Idiom

```go
logger := logrus.New()
logger.SetFormatter(&dispatcher.IBMFormatter{Width: 80})

// Standard logrus usage with IBM formatting
logger.WithFields(logrus.Fields{
    "id":       "APP001",
    "severity": "INFO",
    "service":  "payment-processor",
}).Info("Application started")
```

## Formatter Options

### IBMFormatter

- `Width` - Box width in characters (default: 80)
- `DisableColors` - Disable ANSI color codes (default: false)
- `TimestampFormat` - Timestamp format (default: RFC3339)

### SimpleIBMFormatter

- `DisableColors` - Disable ANSI color codes (default: false)
- `TimestampFormat` - Timestamp format (default: RFC3339)

## Example Output

### Full Box Style (IBMFormatter)

```
╭──────────────────────────────────────────────────────────────────────────────╮
│ [2025-11-13T13:05:00Z] SRV002 (ERROR): Failed to bind to port 8080          │
│ Failed to bind to port {port}                                                │
│                                                                               │
│  port=8080                                                                    │
│  error=address already in use                                                │
│                                                                               │
│ Help: Cause: Port already in use or insufficient privileges. Recovery:      │
│       Stop the conflicting process or choose another port.                   │
╰──────────────────────────────────────────────────────────────────────────────╯
```

### Simple Style (SimpleIBMFormatter)

```
[2025-11-13T13:05:00Z] SRV002 (ERROR): Failed to bind to port 8080
    port=8080
    error=address already in use
    Help: Cause: Port already in use. Recovery: Stop conflicting process.
```

## Integration with Existing Loggers

The formatters work seamlessly with existing Logrus-based applications. Simply replace your formatter:

```go
// Before
logger.SetFormatter(&logrus.JSONFormatter{})

// After - IBM style!
logger.SetFormatter(&dispatcher.IBMFormatter{Width: 80})
```

All your existing log calls will now use IBM-style formatting with beautiful box borders and colors!
