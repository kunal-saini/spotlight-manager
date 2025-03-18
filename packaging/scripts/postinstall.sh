#!/bin/sh
# packaging/scripts/postinstall.sh

set -e

# Reload systemd
systemctl daemon-reload

# Enable and start the service
systemctl enable spotlightd.service
systemctl start spotlightd.service || true

# Create desktop entry for autostart (alternative method to systemd)
mkdir -p /etc/xdg/autostart
cat > /etc/xdg/autostart/spotlightd.desktop << EOF
[Desktop Entry]
Type=Application
Name=Spotlight Wallpaper Manager
Comment=Fetches and sets Windows Spotlight wallpapers
Exec=/usr/bin/spotlightd
Terminal=false
Categories=Utility;
X-GNOME-Autostart-enabled=true
EOF

echo "Spotlight Manager has been installed successfully!"
echo "The service is now running in the background."
echo "Look for the tray icon in your system tray."

exit 0