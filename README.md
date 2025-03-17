## Ubuntu Spotlight Wallpaper Manager

A lightweight application that brings Microsoft's Spotlight wallpapers to Ubuntu 24.04 and other Linux distributions.

### Features

- Automatically fetches high-quality wallpapers from Microsoft's Spotlight API
- Updates wallpaper daily or on demand
- Runs quietly in the system tray
- Preserves metadata and image information
- Configurable refresh interval and storage options
- Installs easily via .deb package

### Installation

#### From .deb package

```bash
sudo dpkg -i spotlight-manager_1.0.0_amd64.deb
```

#### From source

```bash
# Clone the repository
git clone https://github.com/yourusername/spotlight-manager.git
cd spotlight-manager

# Build the application
go build -o spotlightd ./cmd/spotlightd

# Install manually
sudo mkdir -p /usr/bin
sudo cp spotlightd /usr/bin/
sudo cp packaging/systemd/spotlightd.service /lib/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable spotlightd
sudo systemctl start spotlightd
```

### Usage

Once installed, the application will:

1. Run automatically at system startup
2. Update wallpapers daily
3. Show a system tray icon for manual control

#### System Tray Options

- **Refresh Now**: Immediately fetch and apply a new wallpaper
- **Quit**: Exit the application

#### Configuration

The application stores its configuration in `~/.config/spotlight-manager/config.json`:

```json
{
  "wallpaper_dir": "/home/username/.config/spotlight-manager/wallpapers",
  "refresh_interval": 24,
  "keep_count": 10,
  "start_with_system": true
}
```

- `wallpaper_dir`: Directory where wallpapers are stored
- `refresh_interval`: Time between wallpaper updates (in hours)
- `keep_count`: Number of wallpapers to keep before cleaning up (0 = keep all)
- `start_with_system`: Whether to start automatically at system boot

### Building the .deb Package

To build a .deb package for distribution:

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Build the package
goreleaser build --snapshot --clean
```

### Development

#### Project Structure

```
spotlight-manager/
├── cmd/
│   └── spotlightd/
│       └── main.go
├── internal/
│   ├── wallpaper/
│   │   ├── fetcher.go
│   │   └── setter.go
│   ├── config/
│   │   └── config.go
│   └── tray/
│       └── tray.go
├── packaging/
│   ├── systemd/
│   │   └── spotlightd.service
│   └── scripts/
│       ├── postinstall.sh
│       └── preremove.sh
├── .goreleaser.yml
└── go.mod
```

### Credits

- Microsoft Spotlight API reference: [ORelio/Spotlight-Downloader](https://github.com/ORelio/Spotlight-Downloader/blob/master/SpotlightAPI.md)
- System tray implementation: [getlantern/systray](https://github.com/getlantern/systray)

### License

MIT License

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.