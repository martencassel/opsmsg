package main

import (
	"context"
	"log"

	"github.com/martencassel/opsmsg/catalog"
	"github.com/martencassel/opsmsg/dispatcher"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load catalogs
	builtin, err := catalog.Load("../../catalog/builtin.yaml")
	if err != nil {
		log.Fatal("Failed to load builtin catalog: ", err)
	}

	// Demo 1: IBM Formatter with box borders (full style)
	println("═══════════════════════════════════════════════════════════════════════════════")
	println("Demo 1: IBM Formatter (Full Box Style)")
	println("═══════════════════════════════════════════════════════════════════════════════")
	println()

	logger1 := logrus.New()
	logger1.SetFormatter(&dispatcher.IBMFormatter{
		Width: 80,
	})

	d1 := dispatcher.NewLogrusDispatcher(logger1)

	// Dispatch various messages
	msg1 := builtin.New("SRV001", map[string]string{"port": "8080"})
	d1.Dispatch(context.Background(), msg1)

	println() // Space between messages

	msg2 := builtin.New("SRV002", map[string]string{
		"port":  "8080",
		"error": "address already in use",
	})
	d1.Dispatch(context.Background(), msg2)

	println()

	msg3 := builtin.New("DEP002", map[string]string{
		"host":  "db.example.com",
		"error": "connection refused",
	})
	d1.Dispatch(context.Background(), msg3)

	// Demo 2: Simple IBM Formatter (no borders)
	println()
	println("═══════════════════════════════════════════════════════════════════════════════")
	println("Demo 2: Simple IBM Formatter (No Borders)")
	println("═══════════════════════════════════════════════════════════════════════════════")
	println()

	logger2 := logrus.New()
	logger2.SetFormatter(&dispatcher.SimpleIBMFormatter{})

	d2 := dispatcher.NewLogrusDispatcher(logger2)

	msg4 := builtin.New("SEC001", map[string]string{
		"ip":       "192.168.1.100",
		"endpoint": "/admin",
	})
	d2.Dispatch(context.Background(), msg4)

	println()

	msg5 := builtin.New("RTE002", map[string]string{
		"endpoint": "/api/v1/reports",
		"duration": "5.2s",
	})
	d2.Dispatch(context.Background(), msg5)

	// Demo 3: Direct Logrus usage with IBM formatter (logger idiom)
	println()
	println("═══════════════════════════════════════════════════════════════════════════════")
	println("Demo 3: Direct Logrus Logger Idiom")
	println("═══════════════════════════════════════════════════════════════════════════════")
	println()

	logger3 := logrus.New()
	logger3.SetFormatter(&dispatcher.IBMFormatter{Width: 80})

	// Use standard logrus idioms directly
	logger3.WithFields(logrus.Fields{
		"id":       "APP001",
		"severity": "INFO",
		"service":  "payment-processor",
		"version":  "v2.1.0",
	}).Info("Application initialized successfully")

	println()

	logger3.WithFields(logrus.Fields{
		"id":          "NET002",
		"severity":    "ERROR",
		"host":        "api.external.com",
		"port":        "443",
		"retry_count": 3,
		"help":        "Cause: Remote host did not respond. Recovery: Verify host availability and network connectivity.",
	}).Error("Connection timeout")

	println()

	// Demo 4: Disable colors
	println("═══════════════════════════════════════════════════════════════════════════════")
	println("Demo 4: Disable Colors")
	println("═══════════════════════════════════════════════════════════════════════════════")
	println()

	logger4 := logrus.New()
	logger4.SetFormatter(&dispatcher.IBMFormatter{
		Width:         80,
		DisableColors: true,
	})

	d4 := dispatcher.NewLogrusDispatcher(logger4)

	msg6 := builtin.New("SRV003", map[string]string{
		"error": "PORT environment variable not set",
	})
	d4.Dispatch(context.Background(), msg6)
}
