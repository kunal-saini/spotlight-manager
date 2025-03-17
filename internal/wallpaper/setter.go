package wallpaper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kunal-saini/spotlight-manager/internal/config"
)

type Manager struct {
	cfg     *config.Config
	logger  *log.Logger
	fetcher *Fetcher
}

func NewManager(cfg *config.Config, logger *log.Logger) *Manager {
	return &Manager{
		cfg:     cfg,
		logger:  logger,
		fetcher: NewFetcher(logger),
	}
}

func (m *Manager) Refresh() error {
	m.logger.Println("Refreshing wallpaper...")

	// Fetch new wallpaper URL and metadata
	imageURL, metadata, err := m.fetcher.FetchSpotlightMetadata()
	if err != nil {
		return fmt.Errorf("failed to fetch metadata: %w", err)
	}

	m.logger.Printf("Found image: %s", imageURL)

	// Create a timestamp for the filename
	timestamp := time.Now().Format("20060102-150405")

	// Extract a safe filename from the URL
	urlParts := strings.Split(imageURL, "/")
	originalFilename := urlParts[len(urlParts)-1]

	// Create image filename
	imagePath := filepath.Join(m.cfg.WallpaperDir, fmt.Sprintf("spotlight-%s-%s", timestamp, originalFilename))

	// Create metadata filename
	metadataPath := filepath.Join(m.cfg.WallpaperDir, fmt.Sprintf("spotlight-%s.txt", timestamp))

	// Download image
	if err := m.fetcher.DownloadImage(imageURL, imagePath); err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	// Save metadata
	if err := os.WriteFile(metadataPath, []byte(metadata), 0644); err != nil {
		m.logger.Printf("Warning: failed to save metadata: %v", err)
		// Continue even if metadata saving fails
	}

	m.logger.Printf("Saved image to: %s", imagePath)

	// Set dark and light wallpapers
	if err := m.setWallpaper(imagePath); err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}

	// Clean up old wallpapers if configured
	if m.cfg.KeepCount > 0 {
		if err := m.cleanupOldWallpapers(); err != nil {
			m.logger.Printf("Warning: failed to clean up old wallpapers: %v", err)
		}
	}

	return nil
}

func (m *Manager) setWallpaper(imagePath string) error {
	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("wallpaper file does not exist: %s", imagePath)
	}

	// Formulate proper URI for gsettings
	uri := fmt.Sprintf("file://%s", imagePath)

	// Set for dark mode
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", uri)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set dark wallpaper: %w", err)
	}

	// Set for light mode
	cmd = exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", uri)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set light wallpaper: %w", err)
	}

	return nil
}

func (m *Manager) cleanupOldWallpapers() error {
	// Get list of files in wallpaper directory
	files, err := os.ReadDir(m.cfg.WallpaperDir)
	if err != nil {
		return fmt.Errorf("failed to read wallpaper directory: %w", err)
	}

	// Filter for wallpaper files
	var wallpapers []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "spotlight-") && strings.HasSuffix(file.Name(), ".jpg") {
			wallpapers = append(wallpapers, file.Name())
		}
	}

	// Sort wallpapers by name (which includes timestamp)
	// This sorts them chronologically since we use timestamp in the filename
	sort.Strings(wallpapers)

	// Delete oldest wallpapers if we have more than KeepCount
	if len(wallpapers) > m.cfg.KeepCount {
		for i := 0; i < len(wallpapers)-m.cfg.KeepCount; i++ {
			fileToDelete := filepath.Join(m.cfg.WallpaperDir, wallpapers[i])
			if err := os.Remove(fileToDelete); err != nil {
				m.logger.Printf("Warning: failed to delete old wallpaper %s: %v", fileToDelete, err)
			} else {
				m.logger.Printf("Deleted old wallpaper: %s", fileToDelete)
			}

			// Also try to delete the corresponding metadata file
			metadataFile := strings.TrimSuffix(fileToDelete, filepath.Ext(fileToDelete)) + ".txt"
			if err := os.Remove(metadataFile); err != nil && !os.IsNotExist(err) {
				m.logger.Printf("Warning: failed to delete metadata file %s: %v", metadataFile, err)
			}
		}
	}

	return nil
}
