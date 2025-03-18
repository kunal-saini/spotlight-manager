#!/bin/sh
# packaging/scripts/preremove.sh

set -e

# Stop and disable the service
systemctl stop spotlightd.service || true
systemctl disable spotlightd.service || true

# Remove autostart desktop entry if it exists
if [ -f /etc/xdg/autostart/spotlightd.desktop ]; then
    rm -f /etc/xdg/autostart/spotlightd.desktop
fi

echo "Spotlight Manager has been successfully removed."
echo "Note: User configuration and downloaded wallpapers in ~/.config/spotlight-manager/ have been preserved."
echo "To completely remove all data, run: rm -rf ~/.config/spotlight-manager/"

exit 0