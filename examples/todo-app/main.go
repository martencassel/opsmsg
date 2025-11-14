package main

import (
	"context"
	"log"

	"github.com/martencassel/opsmsg/catalog"
	"github.com/martencassel/opsmsg/dispatcher"
	"github.com/sirupsen/logrus"
)

func main() {
	// Setup dispatcher with Logrus
	logger := logrus.New()
	d := dispatcher.NewLogrusDispatcher(logger)

	// Load catalogs
	builtin, err := catalog.Load("../../catalog/builtin.yaml")
	if err != nil {
		log.Fatal(err)
	}
	custom, err := catalog.Load("catalog/custom.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Merge catalogs
	merged := catalog.Merge(builtin, custom)

	// Start server
	StartServer(context.Background(), d, merged)
}
