package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	WallpaperDir    string `json:"wallpaper_dir"`
	RefreshInterval int    `json:"refresh_interval"` // in hours
	KeepCount       int    `json:"keep_count"`       // number of wallpapers to keep (0 = keep all)
	StartWithSystem bool   `json:"start_with_system"`
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".config", "spotlight-manager")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := &Config{
			WallpaperDir:    filepath.Join(configDir, "wallpapers"),
			RefreshInterval: 24,
			KeepCount:       10,
			StartWithSystem: true,
		}

		if err := os.MkdirAll(cfg.WallpaperDir, 0755); err != nil {
			return nil, err
		}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	configDir := filepath.Dir(filepath.Join(c.WallpaperDir, ".."))
	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
