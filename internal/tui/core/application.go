package core

import (
	"context"

	"github.com/rivo/tview"
)

// Application represents the main TUI application
type Application struct {
	app    *tview.Application
	router *Router
	theme  Theme
	layout *LayoutManager
	ctx    context.Context
	cancel context.CancelFunc
}

// NewApplication creates a new Application instance
func NewApplication(theme Theme) *Application {
	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		app:    tview.NewApplication(),
		router: NewRouter(),
		theme:  theme,
		layout: NewLayoutManager(),
		ctx:    ctx,
		cancel: cancel,
	}
}

// GetApplication returns the tview application
func (a *Application) GetApplication() *tview.Application {
	return a.app
}

// GetRouter returns the router
func (a *Application) GetRouter() *Router {
	return a.router
}

// GetTheme returns the theme
func (a *Application) GetTheme() Theme {
	return a.theme
}

// GetLayoutManager returns the layout manager
func (a *Application) GetLayoutManager() *LayoutManager {
	return a.layout
}

// GetContext returns the context
func (a *Application) GetContext() context.Context {
	return a.ctx
}

// Cancel cancels the context
func (a *Application) Cancel() {
	a.cancel()
}

// Run starts the application
func (a *Application) Run() error {
	defer a.cancel()
	return a.app.Run()
}

// Stop stops the application
func (a *Application) Stop() {
	a.cancel()
	a.app.Stop()
}
