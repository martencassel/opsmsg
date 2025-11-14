package dispatcher

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

// IBMFormatter formats log entries in IBM-style textual format with box borders
type IBMFormatter struct {
	// DisableColors disables ANSI color output
	DisableColors bool
	// Width sets the box width (default: 80)
	Width int
	// TimestampFormat sets the timestamp format (default: RFC3339)
	TimestampFormat string
}

// Format renders an Entry in IBM-style format with box borders
func (f *IBMFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	width := f.Width
	if width == 0 {
		width = 80
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var b strings.Builder

	// Extract message fields
	msgID, _ := entry.Data["id"].(string)
	if msgID == "" {
		msgID = "UNKNOWN"
	}

	severity := f.getSeverityString(entry)
	timestamp := entry.Time.Format(timestampFormat)
	text := entry.Message

	// Get colors (or empty strings if disabled)
	orange := f.color(colorOrange)
	severityColor := f.getSeverityColor(entry.Level)
	gray := f.color(colorGray)
	bright := f.color(colorBright)
	dim := f.color(colorDim)
	yellow := f.color(colorYellow)
	reset := f.color(colorReset)

	// Top border
	b.WriteString(orange)
	b.WriteString(boxTopLeft)
	b.WriteString(strings.Repeat(boxHorizontal, width-2))
	b.WriteString(boxTopRight)
	b.WriteString(reset)
	b.WriteString("\n")

	// Timestamp, ID and Severity line
	b.WriteString(orange)
	b.WriteString(boxVertical)
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(gray)
	b.WriteString("[")
	b.WriteString(timestamp)
	b.WriteString("]")
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(bright)
	b.WriteString(orange)
	b.WriteString(msgID)
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(severityColor)
	b.WriteString("(")
	b.WriteString(severity)
	b.WriteString(")")
	b.WriteString(reset)

	// Pad to width
	currentLen := len(timestamp) + len(msgID) + len(severity) + 7
	padding := width - currentLen - 4
	if padding < 0 {
		padding = 0
	}
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(orange)
	b.WriteString(boxVertical)
	b.WriteString(reset)
	b.WriteString("\n")

	// Message text - wrap if needed
	textLines := f.wrapText(text, width-4)
	for _, line := range textLines {
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString(" ")
		b.WriteString(bright)
		b.WriteString(line)
		b.WriteString(reset)

		textLen := len(line)
		padding = width - textLen - 4
		if padding < 0 {
			padding = 0
		}
		b.WriteString(strings.Repeat(" ", padding))
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString("\n")
	}

	// Context fields (excluding internal fields)
	contextFields := make(map[string]interface{})
	for k, v := range entry.Data {
		if k != "id" && k != "severity" && k != "timestamp" && k != "help" {
			contextFields[k] = v
		}
	}

	if len(contextFields) > 0 {
		// Empty line separator
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString(strings.Repeat(" ", width-2))
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString("\n")

		// Print context fields with wrapping
		for k, v := range contextFields {
			line := fmt.Sprintf("    %s=%v", k, v)
			// Wrap if line is too long
			if len(line) > width-4 {
				// Wrap the line
				wrappedLines := f.wrapText(line, width-6)
				for i, wrappedLine := range wrappedLines {
					b.WriteString(orange)
					b.WriteString(boxVertical)
					b.WriteString(reset)
					if i == 0 {
						b.WriteString(" ")
					} else {
						b.WriteString("      ") // indent continuation
					}
					b.WriteString(dim)
					b.WriteString(wrappedLine)
					b.WriteString(reset)

					lineLen := len(wrappedLine) + 1
					if i > 0 {
						lineLen = len(wrappedLine) + 6
					}
					padding = width - lineLen - 2
					if padding < 0 {
						padding = 0
					}
					b.WriteString(strings.Repeat(" ", padding))
					b.WriteString(orange)
					b.WriteString(boxVertical)
					b.WriteString(reset)
					b.WriteString("\n")
				}
			} else {
				b.WriteString(orange)
				b.WriteString(boxVertical)
				b.WriteString(reset)
				b.WriteString(" ")
				b.WriteString(dim)
				b.WriteString(line)
				b.WriteString(reset)

				lineLen := len(line) + 1
				padding = width - lineLen - 2
				if padding < 0 {
					padding = 0
				}
				b.WriteString(strings.Repeat(" ", padding))
				b.WriteString(orange)
				b.WriteString(boxVertical)
				b.WriteString(reset)
				b.WriteString("\n")
			}
		}
	}

	// Help text if present
	if helpText, ok := entry.Data["help"].(string); ok && helpText != "" {
		// Empty line separator
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString(strings.Repeat(" ", width-2))
		b.WriteString(orange)
		b.WriteString(boxVertical)
		b.WriteString(reset)
		b.WriteString("\n")

		// Help section
		helpLines := f.wrapText(helpText, width-10)
		for i, line := range helpLines {
			b.WriteString(orange)
			b.WriteString(boxVertical)
			b.WriteString(reset)
			if i == 0 {
				b.WriteString(" ")
				b.WriteString(yellow)
				b.WriteString("Help:")
				b.WriteString(reset)
				b.WriteString(" ")
			} else {
				b.WriteString("       ")
			}
			b.WriteString(dim)
			b.WriteString(line)
			b.WriteString(reset)

			var totalLen int
			if i == 0 {
				totalLen = 1 + 6 + 1 + len(line) + 1 + 1
			} else {
				totalLen = 1 + 7 + len(line) + 1 + 1
			}
			padding = width - totalLen
			if padding < 0 {
				padding = 0
			}
			b.WriteString(strings.Repeat(" ", padding))
			b.WriteString(orange)
			b.WriteString(boxVertical)
			b.WriteString(reset)
			b.WriteString("\n")
		}
	}

	// Bottom border
	b.WriteString(orange)
	b.WriteString(boxBottomLeft)
	b.WriteString(strings.Repeat(boxHorizontal, width-2))
	b.WriteString(boxBottomRight)
	b.WriteString(reset)
	b.WriteString("\n")

	return []byte(b.String()), nil
}

func (f *IBMFormatter) color(code string) string {
	if f.DisableColors {
		return ""
	}
	return code
}

func (f *IBMFormatter) getSeverityString(entry *logrus.Entry) string {
	if sev, ok := entry.Data["severity"].(string); ok {
		return sev
	}
	// Fallback to logrus level
	return strings.ToUpper(entry.Level.String())
}

func (f *IBMFormatter) getSeverityColor(level logrus.Level) string {
	if f.DisableColors {
		return ""
	}
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel:
		return colorRed
	case logrus.ErrorLevel:
		return colorRed
	case logrus.WarnLevel:
		return colorYellow
	case logrus.InfoLevel:
		return colorCyan
	default:
		return colorOrange
	}
}

func (f *IBMFormatter) wrapText(text string, width int) []string {
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

// SimpleIBMFormatter provides a simpler IBM-style format without box borders
type SimpleIBMFormatter struct {
	// DisableColors disables ANSI color output
	DisableColors bool
	// TimestampFormat sets the timestamp format (default: RFC3339)
	TimestampFormat string
}

// Format renders an Entry in simple IBM-style format
func (f *SimpleIBMFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	var b strings.Builder

	// Extract message fields
	msgID, _ := entry.Data["id"].(string)
	if msgID == "" {
		msgID = "UNKNOWN"
	}

	severity := f.getSeverityString(entry)
	timestamp := entry.Time.Format(timestampFormat)

	// Get colors
	orange := f.color(colorOrange)
	severityColor := f.getSeverityColor(entry.Level)
	gray := f.color(colorGray)
	bright := f.color(colorBright)
	dim := f.color(colorDim)
	yellow := f.color(colorYellow)
	reset := f.color(colorReset)

	// Header line with colors
	b.WriteString(gray)
	b.WriteString("[")
	b.WriteString(timestamp)
	b.WriteString("]")
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(bright)
	b.WriteString(orange)
	b.WriteString(msgID)
	b.WriteString(reset)
	b.WriteString(" ")
	b.WriteString(severityColor)
	b.WriteString("(")
	b.WriteString(severity)
	b.WriteString(")")
	b.WriteString(reset)
	b.WriteString(": ")
	b.WriteString(bright)
	b.WriteString(entry.Message)
	b.WriteString(reset)
	b.WriteString("\n")

	// Context fields
	for k, v := range entry.Data {
		if k != "id" && k != "severity" && k != "timestamp" && k != "help" {
			b.WriteString("    ")
			b.WriteString(dim)
			fmt.Fprintf(&b, "%s=%v", k, v)
			b.WriteString(reset)
			b.WriteString("\n")
		}
	}

	// Help text
	if helpText, ok := entry.Data["help"].(string); ok && helpText != "" {
		b.WriteString("    ")
		b.WriteString(yellow)
		b.WriteString("Help:")
		b.WriteString(reset)
		b.WriteString(" ")
		b.WriteString(dim)
		b.WriteString(helpText)
		b.WriteString(reset)
		b.WriteString("\n")
	}

	return []byte(b.String()), nil
}

func (f *SimpleIBMFormatter) color(code string) string {
	if f.DisableColors {
		return ""
	}
	return code
}

func (f *SimpleIBMFormatter) getSeverityString(entry *logrus.Entry) string {
	if sev, ok := entry.Data["severity"].(string); ok {
		return sev
	}
	return strings.ToUpper(entry.Level.String())
}

func (f *SimpleIBMFormatter) getSeverityColor(level logrus.Level) string {
	if f.DisableColors {
		return ""
	}
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel:
		return colorRed
	case logrus.ErrorLevel:
		return colorRed
	case logrus.WarnLevel:
		return colorYellow
	case logrus.InfoLevel:
		return colorCyan
	default:
		return colorOrange
	}
}
