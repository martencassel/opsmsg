package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/martencassel/opsmsg/catalog"
	"github.com/martencassel/opsmsg/message"
)

const (
	// ANSI color codes
	colorReset  = "\033[0m"
	colorOrange = "\033[38;5;208m"
	colorYellow = "\033[38;5;226m"
	colorRed    = "\033[38;5;196m"
	colorCyan   = "\033[38;5;51m"
	colorGray   = "\033[38;5;240m"
	colorBright = "\033[1m"
	colorDim    = "\033[2m"

	// Box drawing characters
	boxTopLeft     = "╭"
	boxTopRight    = "╮"
	boxBottomLeft  = "╰"
	boxBottomRight = "╯"
	boxHorizontal  = "─"
	boxVertical    = "│"
)

type MessageViewer struct {
	messages []message.Message
	current  int
}

func NewMessageViewer() *MessageViewer {
	return &MessageViewer{
		messages: []message.Message{},
		current:  0,
	}
}

func (v *MessageViewer) AddMessage(msg message.Message) {
	v.messages = append(v.messages, msg)
}

func (v *MessageViewer) getSeverityColor(severity message.Severity) string {
	switch severity {
	case message.Critical:
		return colorRed
	case message.Error:
		return colorRed
	case message.Warn:
		return colorYellow
	case message.Info:
		return colorCyan
	default:
		return colorOrange
	}
}

func (v *MessageViewer) renderMessage(msg message.Message) string {
	var sb strings.Builder

	// Top border
	width := 80
	sb.WriteString(colorOrange)
	sb.WriteString(boxTopLeft)
	sb.WriteString(strings.Repeat(boxHorizontal, width-2))
	sb.WriteString(boxTopRight)
	sb.WriteString(colorReset)
	sb.WriteString("\n")

	// Timestamp, ID and Severity line
	timestamp := msg.Timestamp.Format("2006-01-02T15:04:05Z")
	severityColor := v.getSeverityColor(msg.Severity)

	sb.WriteString(colorOrange)
	sb.WriteString(boxVertical)
	sb.WriteString(colorReset)
	sb.WriteString(" ")
	sb.WriteString(colorGray)
	sb.WriteString("[")
	sb.WriteString(timestamp)
	sb.WriteString("]")
	sb.WriteString(colorReset)
	sb.WriteString(" ")
	sb.WriteString(colorBright)
	sb.WriteString(colorOrange)
	sb.WriteString(msg.ID)
	sb.WriteString(colorReset)
	sb.WriteString(" ")
	sb.WriteString(severityColor)
	sb.WriteString("(")
	sb.WriteString(string(msg.Severity))
	sb.WriteString(")")
	sb.WriteString(colorReset)

	// Pad to width
	currentLen := len(timestamp) + len(msg.ID) + len(msg.Severity) + 7 // brackets, parens, spaces
	padding := width - currentLen - 4
	if padding < 0 {
		padding = 0
	}
	sb.WriteString(strings.Repeat(" ", padding))
	sb.WriteString(colorOrange)
	sb.WriteString(boxVertical)
	sb.WriteString(colorReset)
	sb.WriteString("\n")

	// Message text - wrap if needed
	text := msg.Text
	for k, v := range msg.Context {
		text = strings.ReplaceAll(text, "{"+k+"}", v)
	}

	textLines := wrapText(text, width-4)
	for _, line := range textLines {
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString(" ")
		sb.WriteString(colorBright)
		sb.WriteString(line)
		sb.WriteString(colorReset)

		textLen := len(line)
		padding = width - textLen - 4
		if padding < 0 {
			padding = 0
		}
		sb.WriteString(strings.Repeat(" ", padding))
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString("\n")
	}

	// Context section (if any)
	if len(msg.Context) > 0 {
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString(strings.Repeat(" ", width-2))
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString("\n")

		for k, v := range msg.Context {
			sb.WriteString(colorOrange)
			sb.WriteString(boxVertical)
			sb.WriteString(colorReset)
			sb.WriteString("    ")
			sb.WriteString(colorDim)
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			sb.WriteString(colorReset)

			contextLen := len(k) + len(v) + 5 // "    ", "=", and spaces
			padding = width - contextLen - 2
			if padding < 0 {
				padding = 0
			}
			sb.WriteString(strings.Repeat(" ", padding))
			sb.WriteString(colorOrange)
			sb.WriteString(boxVertical)
			sb.WriteString(colorReset)
			sb.WriteString("\n")
		}
	}

	// Help section
	if msg.Help != "" {
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString(strings.Repeat(" ", width-2))
		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString("\n")

		sb.WriteString(colorOrange)
		sb.WriteString(boxVertical)
		sb.WriteString(colorReset)
		sb.WriteString(" ")
		sb.WriteString(colorYellow)
		sb.WriteString("Help:")
		sb.WriteString(colorReset)
		sb.WriteString(" ")
		sb.WriteString(colorDim)

		// Word wrap help text
		helpLines := wrapText(msg.Help, width-10)
		for i, line := range helpLines {
			if i > 0 {
				sb.WriteString(colorOrange)
				sb.WriteString(boxVertical)
				sb.WriteString(colorReset)
				sb.WriteString("       ")
			}
			sb.WriteString(line)
			sb.WriteString(colorReset)

			lineLen := len(line)
			var totalLen int
			if i == 0 {
				totalLen = 1 + 6 + 1 + lineLen + 1 + 1 // │ + "Help: " + " " + line + " " + │
			} else {
				totalLen = 1 + 7 + lineLen + 1 + 1 // │ + "       " + line + " " + │
			}
			padding = width - totalLen
			if padding < 0 {
				padding = 0
			}
			sb.WriteString(strings.Repeat(" ", padding))
			sb.WriteString(colorOrange)
			sb.WriteString(boxVertical)
			sb.WriteString(colorReset)
			sb.WriteString("\n")
		}
	}

	// Bottom border
	sb.WriteString(colorOrange)
	sb.WriteString(boxBottomLeft)
	sb.WriteString(strings.Repeat(boxHorizontal, width-2))
	sb.WriteString(boxBottomRight)
	sb.WriteString(colorReset)
	sb.WriteString("\n")

	return sb.String()
}

func wrapText(text string, width int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+len(word)+1 <= width {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

func (v *MessageViewer) Show() {
	if len(v.messages) == 0 {
		fmt.Println("No messages to display")
		return
	}

	// Clear screen
	fmt.Print("\033[2J\033[H")

	for {
		// Clear and show header
		fmt.Print("\033[H")
		fmt.Printf("\n%s%s╔═══════════════════════════════════════════════════════════════════════════════╗%s\n", colorBright, colorOrange, colorReset)
		fmt.Printf("%s%s║%s                     %s%sOPMSG Message Viewer%s                                  %s%s║%s\n", colorBright, colorOrange, colorReset, colorBright, colorOrange, colorReset, colorBright, colorOrange, colorReset)
		fmt.Printf("%s%s╚═══════════════════════════════════════════════════════════════════════════════╝%s\n\n", colorBright, colorOrange, colorReset)

		// Show current message
		fmt.Println(v.renderMessage(v.messages[v.current]))

		// Show navigation info
		fmt.Printf("\n%s%sMessage %d of %d%s\n", colorDim, colorGray, v.current+1, len(v.messages), colorReset)
		fmt.Printf("%s%sPress [n]ext, [p]revious, or [q]uit: %s", colorBright, colorOrange, colorReset)

		// Read single character
		var input string
		fmt.Scanln(&input)

		if len(input) == 0 {
			continue
		}

		switch strings.ToLower(string(input[0])) {
		case "n":
			if v.current < len(v.messages)-1 {
				v.current++
			}
		case "p":
			if v.current > 0 {
				v.current--
			}
		case "q":
			fmt.Println("\n" + colorOrange + "Goodbye!" + colorReset)
			return
		}

		// Clear screen for next iteration
		fmt.Print("\033[2J")
	}
}

func main() {
	// Load catalog
	cat, err := catalog.Load("../../catalog/builtin.yaml")
	if err != nil {
		log.Fatal("Failed to load builtin catalog: ", err)
	}

	customCat, err := catalog.Load("catalog/messages.yaml")
	if err != nil {
		log.Fatal("Failed to load custom catalog: ", err)
	}

	merged := catalog.Merge(cat, customCat)

	// Create viewer
	viewer := NewMessageViewer()

	// Generate sample messages with various contexts
	ctx := context.Background()
	_ = ctx

	// Sample messages
	messages := []struct {
		id      string
		context map[string]string
	}{
		{"SRV001", map[string]string{"port": "8080"}},
		{"SRV002", map[string]string{"port": "8080", "error": "address already in use"}},
		{"DEP002", map[string]string{"host": "db.example.com", "error": "connection refused"}},
		{"API002", map[string]string{"client_id": "client-12345", "current": "950", "limit": "1000"}},
		{"APP003", map[string]string{"job_id": "batch-2024-001", "error": "invalid data format"}},
		{"NET002", map[string]string{"host": "api.external.com", "port": "443"}},
		{"SEC001", map[string]string{"ip": "192.168.1.100", "endpoint": "/admin"}},
		{"APP004", map[string]string{"table_name": "users", "records_affected": "1523"}},
		{"RTE002", map[string]string{"endpoint": "/api/v1/reports", "duration": "5.2s"}},
		{"SRV003", map[string]string{"error": "PORT environment variable not set"}},
	}

	// Set different timestamps for variety
	baseTime := time.Now().Add(-1 * time.Hour)
	for i, msgDef := range messages {
		msg := merged.New(msgDef.id, msgDef.context)
		msg.Timestamp = baseTime.Add(time.Duration(i*5) * time.Minute)
		viewer.AddMessage(msg)
	}

	// Check if terminal supports colors
	if os.Getenv("TERM") == "" {
		fmt.Println("Warning: Terminal type not detected, colors may not display correctly")
	}

	// Show the viewer
	viewer.Show()
}
