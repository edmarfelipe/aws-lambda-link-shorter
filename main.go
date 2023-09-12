package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/apex/gateway/v2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/edmarfelipe/aws-lambda/handler/createlink"
	"github.com/edmarfelipe/aws-lambda/handler/redirectlink"
	"github.com/edmarfelipe/aws-lambda/storage"
)

func main() {
	if err := run(); err != nil {
		slog.Info("Failed to start: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	linkStorage := storage.NewLinkStorage(cfg)
	http.Handle("/", redirectlink.NewHandler(linkStorage))
	http.Handle("/link", createlink.NewHandler(linkStorage))
	return gateway.ListenAndServe(":3000", nil)
}
