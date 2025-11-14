package main

import (
	"context"
	"net/http"
	"os"

	"github.com/martencassel/opsmsg/catalog"
	"github.com/martencassel/opsmsg/dispatcher"
)

func StartServer(ctx context.Context, d dispatcher.Dispatcher, c catalog.Catalog) {
	port := os.Getenv("PORT")
	if port == "" {
		msg := c.New("SRV003", map[string]string{"error": "PORT env var missing"})
		d.Dispatch(ctx, msg)
		os.Exit(1)
	}

	addr := ":" + port
	srv := &http.Server{Addr: addr, Handler: routes(d, c)}

	msg := c.New("SRV001", map[string]string{"port": port})
	d.Dispatch(ctx, msg)

	if err := srv.ListenAndServe(); err != nil {
		msg := c.New("SRV002", map[string]string{"port": port, "error": err.Error()})
		d.Dispatch(ctx, msg)
		os.Exit(1)
	}
}
