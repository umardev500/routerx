package routerx

import (
	"net/http"
)

type Group struct {
	prefix string
	app    *App
}

func NewGroup(prefix string, app *App) *Group {
	return &Group{
		prefix: prefix,
		app:    app,
	}
}

// Delete implements Router
func (g *Group) Delete(path string, handlers ...Handler) {
	path = g.prefix + path
	g.app.handle(http.MethodDelete, path, handlers...)
}

// Get implements Router
func (g *Group) Get(path string, handlers ...Handler) {
	path = g.prefix + path
	g.app.handle(http.MethodGet, path, handlers...)
}

// Grou implements Router
func (g *Group) Group(prefix string) Router {
	return NewGroup(g.prefix+prefix, g.app)
}

// Patch implements Router
func (g *Group) Patch(path string, handlers ...Handler) {
	path = g.prefix + path
	g.app.handle(http.MethodPatch, path, handlers...)
}

// Post implements Router
func (g *Group) Post(path string, handlers ...Handler) {
	path = g.prefix + path
	g.app.handle(http.MethodPost, path, handlers...)
}

// Put implements Router
func (g *Group) Put(path string, handlers ...Handler) {
	path = g.prefix + path
	g.app.handle(http.MethodPut, path, handlers...)
}

// Use implements Router
func (g *Group) Use(handlers ...Handler) {
	g.app.middlewares[g.prefix] = append(g.app.middlewares[g.prefix], handlers...)
}
