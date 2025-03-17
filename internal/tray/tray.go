package tray

import (
	"context"
	"log"

	"github.com/getlantern/systray"
	"github.com/kunal-saini/spotlight-manager/internal/wallpaper"
)

type TrayIcon struct {
	ctx     context.Context
	manager *wallpaper.Manager
	logger  *log.Logger
	quitCh  chan struct{}
}

func New(ctx context.Context, manager *wallpaper.Manager, logger *log.Logger) (*TrayIcon, error) {
	t := &TrayIcon{
		ctx:     ctx,
		manager: manager,
		logger:  logger,
		quitCh:  make(chan struct{}),
	}

	go systray.Run(t.onReady, t.onExit)

	return t, nil
}

func (t *TrayIcon) onReady() {
	systray.SetTitle("Spotlight")
	systray.SetTooltip("Ubuntu Spotlight Wallpaper Manager")

	refresh := systray.AddMenuItem("Refresh Now", "Refresh wallpaper immediately")
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			case <-refresh.ClickedCh:
				if err := t.manager.Refresh(); err != nil {
					t.logger.Printf("Failed to refresh wallpaper: %v", err)
				}
			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (t *TrayIcon) onExit() {
	close(t.quitCh)
}

func (t *TrayIcon) Quit() {
	systray.Quit()
	<-t.quitCh
}
