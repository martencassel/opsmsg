package dispatcher

import (
	"context"

	"github.com/martencassel/opsmsg/message"
	"github.com/sirupsen/logrus"
)

type LogrusDispatcher struct {
	logger *logrus.Logger
}

func NewLogrusDispatcher(logger *logrus.Logger) *LogrusDispatcher {
	return &LogrusDispatcher{logger: logger}
}

func (d *LogrusDispatcher) Dispatch(ctx context.Context, msg message.Message) error {
	fields := logrus.Fields{
		"id":       msg.ID,
		"severity": string(msg.Severity),
	}

	// Add context fields
	for k, v := range msg.Context {
		fields[k] = v
	}

	// Add help text if present
	if msg.Help != "" {
		fields["help"] = msg.Help
	}

	entry := d.logger.WithFields(fields)

	switch msg.Severity {
	case message.Info:
		entry.Info(msg.Text)
	case message.Warn:
		entry.Warn(msg.Text)
	case message.Error:
		entry.Error(msg.Text)
	case message.Critical:
		entry.Fatal(msg.Text)
	default:
		entry.Info(msg.Text)
	}

	return nil
}
