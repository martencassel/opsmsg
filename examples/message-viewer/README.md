# Message Viewer Example

A simple interactive terminal application that displays operational messages with pagination and IBM-style textual UI formatting.

## Features

- 🎨 Beautiful IBM-inspired text UI with box-drawing characters
- 🟠 Orange color scheme with context-aware severity colors
- 📄 Paginated message browsing (one message at a time)
- 🔍 Full message details including context variables and help text
- ⌨️  Simple keyboard navigation

## Running

From this directory:

```bash
go run main.go
```

## Navigation

- `n` - Next message
- `p` - Previous message
- `q` - Quit

## Example Output

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                     OPMSG Message Viewer                                      ║
╚═══════════════════════════════════════════════════════════════════════════════╝

╭──────────────────────────────────────────────────────────────────────────────╮
│ [2025-11-13T13:10:00Z] SRV002 (ERROR): Failed to bind to port 8080          │
│ Failed to bind to port {port}                                                │
│                                                                               │
│    port=8080                                                                 │
│    error=address already in use                                              │
│                                                                               │
│ Help: Cause: Port already in use or insufficient privileges. Recovery:      │
│       Stop the conflicting process or choose another port.                   │
╰──────────────────────────────────────────────────────────────────────────────╯

Message 2 of 10
Press [n]ext, [p]revious, or [q]uit:
```

## What It Demonstrates

This example shows how to:

1. Load and merge multiple message catalogs
2. Create messages with context variables
3. Render messages in a user-friendly format
4. Build interactive CLI applications with the opsmsg library
5. Format operational messages for human consumption

## Message Catalogs

The viewer loads two catalogs:

- `../../catalog/builtin.yaml` - Built-in system messages
- `catalog/messages.yaml` - Custom application messages

Messages are displayed with:
- Timestamp
- Message ID
- Severity level (with color coding)
- Message text (with expanded context variables)
- Context key-value pairs
- Help text with cause and recovery information
