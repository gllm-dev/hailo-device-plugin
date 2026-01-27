package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gllm-dev/hailo-device-plugin/internal/config"
	"github.com/gllm-dev/hailo-device-plugin/internal/infra/detector"
	"github.com/gllm-dev/hailo-device-plugin/internal/plugin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()

	logger.Info("starting hailo device plugin",
		slog.String("resource_name", cfg.ResourceName),
		slog.String("architecture", cfg.Architecture),
	)

	det := detector.New(cfg)
	p := plugin.New(cfg, det, logger)

	if err := p.Start(); err != nil {
		logger.Error("failed to start plugin", slog.Any("error", err))
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("received shutdown signal", slog.String("signal", sig.String()))

	p.Stop()
	logger.Info("hailo device plugin stopped")
}
