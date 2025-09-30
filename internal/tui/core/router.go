package core

import (
	"fmt"

	"github.com/rivo/tview"
	"slv.sh/slv/internal/tui/interfaces"
)

// Router handles page navigation and routing
type Router struct {
	pages         *tview.Pages
	pageStack     []string
	currentPage   string
	pageRegistry  map[string]interfaces.Page        // Registry of Page interfaces (legacy)
	pageFactories map[string]interfaces.PageFactory // Registry of Page factories
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		pages:         tview.NewPages(),
		pageStack:     make([]string, 0),
		currentPage:   "",
		pageRegistry:  make(map[string]interfaces.Page),
		pageFactories: make(map[string]interfaces.PageFactory),
	}
}

// GetPages returns the pages container
func (r *Router) GetPages() *tview.Pages {
	return r.pages
}

// GetPageStack returns the page stack
func (r *Router) GetPageStack() []string {
	return r.pageStack
}

// GetCurrentPage returns the current page name
func (r *Router) GetCurrentPage() string {
	return r.currentPage
}

// SetCurrentPage sets the current page
func (r *Router) SetCurrentPage(pageName string) {
	r.currentPage = pageName
}

// PushPage pushes a page onto the stack
func (r *Router) PushPage(pageName string) {
	r.pageStack = append(r.pageStack, pageName)
}

// PopPage pops a page from the stack
func (r *Router) PopPage() string {
	if len(r.pageStack) == 0 {
		return ""
	}

	lastIndex := len(r.pageStack) - 1
	pageName := r.pageStack[lastIndex]
	r.pageStack = r.pageStack[:lastIndex]
	return pageName
}

// ClearStack clears the page stack
func (r *Router) ClearStack() {
	r.pageStack = make([]string, 0)
}

// AddPage adds a page to the router
func (r *Router) AddPage(name string, page tview.Primitive, resize bool, visible bool) {
	r.pages.AddPage(name, page, resize, visible)
}

// RemovePage removes a page from the router
func (r *Router) RemovePage(name string) {
	r.pages.RemovePage(name)
}

// HasPage checks if a page exists
func (r *Router) HasPage(name string) bool {
	return r.pages.HasPage(name)
}

// ===== INFRASTRUCTURE METHODS (to avoid duplication in Navigation) =====

// AddPageToMainComponent adds a page to the main content component
func (r *Router) AddPageToMainComponent(name string, page tview.Primitive, components interfaces.ComponentManagerInterface) {
	components.GetMainContentPages().AddPage(name, page, true, false)
}

// NavigateToPage navigates to a page with full stack management and component switching
func (r *Router) NavigateToPage(name string, components interfaces.ComponentManagerInterface, replace bool) {
	if replace {
		// Replace mode: don't push to stack, just switch
		r.currentPage = name
		components.GetMainContentPages().SwitchToPage(name)
	} else {
		// Normal mode: push current page to stack if it exists
		if r.currentPage != "" {
			r.pageStack = append(r.pageStack, r.currentPage)
		}

		// Set new current page
		r.currentPage = name

		// Switch to the page in components (this is the main UI switching)
		components.GetMainContentPages().SwitchToPage(name)
	}
}

// GoBackWithComponents goes back using the stack and updates components
func (r *Router) GoBackWithComponents(components interfaces.ComponentManagerInterface) error {
	if len(r.pageStack) == 0 {
		return fmt.Errorf("no pages in stack to go back to")
	}

	// Pop previous page from stack
	previousPage := r.pageStack[len(r.pageStack)-1]
	r.pageStack = r.pageStack[:len(r.pageStack)-1]

	// Set current page
	r.currentPage = previousPage

	// Switch to previous page in components
	components.GetMainContentPages().SwitchToPage(previousPage)

	return nil
}

// RegisterPage registers a Page interface with the router
func (r *Router) RegisterPage(name string, page interfaces.Page) {
	r.pageRegistry[name] = page
}

// GetRegisteredPage gets a registered Page interface
func (r *Router) GetRegisteredPage(name string) interfaces.Page {
	return r.pageRegistry[name]
}

// HasRegisteredPage checks if a page is registered
func (r *Router) HasRegisteredPage(name string) bool {
	_, exists := r.pageRegistry[name]
	return exists
}

// GetRegisteredPageNames returns all registered page names
func (r *Router) GetRegisteredPageNames() []string {
	names := make([]string, 0, len(r.pageRegistry))
	for name := range r.pageRegistry {
		names = append(names, name)
	}
	return names
}

// RegisterPageFactory registers a Page factory with the router
func (r *Router) RegisterPageFactory(name string, factory interfaces.PageFactory) {
	r.pageFactories[name] = factory
}

// CreatePage creates a new page instance using the registered factory
func (r *Router) CreatePage(tui interfaces.TUIInterface, name string, params ...interface{}) interfaces.Page {
	if factory, exists := r.pageFactories[name]; exists {
		// Prepend TUI to params for the factory
		allParams := append([]interface{}{tui}, params...)
		return factory.CreatePage(allParams...)
	}
	return nil
}

// HasPageFactory checks if a page factory is registered
func (r *Router) HasPageFactory(name string) bool {
	_, exists := r.pageFactories[name]
	return exists
}

// GetPageFactoryNames returns all registered page factory names
func (r *Router) GetPageFactoryNames() []string {
	names := make([]string, 0, len(r.pageFactories))
	for name := range r.pageFactories {
		names = append(names, name)
	}
	return names
}
