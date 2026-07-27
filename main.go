package main

import (
	"log/slog"
	"net/http"
	"os"

	kmf1 "github.com/KohanMathers/KmF1/go"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	client, err := kmf1.New("cache.sqlite", kmf1.WithLogger(logger))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	addr := ":8080"
	logger.Info("starting http server", "addr", addr)
	if err := http.ListenAndServe(addr, newMux(client, logger)); err != nil {
		logger.Error("http server stopped", "error", err)
		panic(err)
	}
}
