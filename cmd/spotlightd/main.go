package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kunal-saini/spotlight-manager/internal/config"
	"github.com/kunal-saini/spotlight-manager/internal/tray"
	"github.com/kunal-saini/spotlight-manager/internal/wallpaper"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[spotlightd] ", log.LstdFlags)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize wallpaper manager
	wpManager := wallpaper.NewManager(cfg, logger)

	// Start system tray
	trayIcon, err := tray.New(ctx, wpManager, logger)
	if err != nil {
		logger.Fatalf("Failed to create tray icon: %v", err)
	}

	// Start wallpaper refresh routine
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := wpManager.Refresh(); err != nil {
					logger.Printf("Failed to refresh wallpaper: %v", err)
				}
			}
		}
	}()

	// Handle system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Println("Shutting down...")
	cancel()
	trayIcon.Quit()
}
