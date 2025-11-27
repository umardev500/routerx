package routerx

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	server      *http.Server
	mux         *http.ServeMux
	middlewares map[string][]Handler
}

const (
	rootMiddlewarePrefix = "/"
)

func New() *App {
	mux := http.NewServeMux()
	return &App{
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		mux:         mux,
		middlewares: map[string][]Handler{},
	}
}

// Listen trigger the server to start
func (a *App) Listen(addr *string) error {
	if addr != nil {
		a.server.Addr = *addr
	}

	return a.server.ListenAndServe()
}

// handle register the handle with given path
func (a *App) handle(method, path string, handlers ...Handler) {
	normalizedPath := NormalizePath(path)
	finalPath := fmt.Sprintf("%s %s", method, normalizedPath)

	middlewares := a.collectMiddlewares(path)
	allHandlers := append(middlewares, handlers...)

	a.mux.HandleFunc(finalPath, func(w http.ResponseWriter, r *http.Request) {
		c := NewCtx(w, r, allHandlers)

		if r.Method != method {
			c.Status(http.StatusMethodNotAllowed).
				JSON(http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		if err := c.Next(); err != nil {
			c.Status(http.StatusInternalServerError).
				JSON(err.Error())
			return
		}
	})
}

// Delete implements Router
func (a *App) Delete(path string, handlers ...Handler) {
	a.handle(http.MethodDelete, path, handlers...)
}

// Get implements Router
func (a *App) Get(path string, handlers ...Handler) {
	a.handle(http.MethodGet, path, handlers...)
}

// Grou implements Router
func (a *App) Group(prefix string) Router {
	return NewGroup(prefix, a)
}

// Patch implements Router
func (a *App) Patch(path string, handlers ...Handler) {
	a.handle(http.MethodPatch, path, handlers...)
}

// Post implements Router
func (a *App) Post(path string, handlers ...Handler) {
	a.handle(http.MethodPost, path, handlers...)
}

// Put implements Router
func (a *App) Put(path string, handlers ...Handler) {
	a.handle(http.MethodPut, path, handlers...)
}

// Use implements Router
func (a *App) Use(handlers ...Handler) {
	a.middlewares[rootMiddlewarePrefix] = append(a.middlewares[rootMiddlewarePrefix], handlers...)
}

// collectMiddlewares returns a slice of middleware handlers that should be applied
// for a given route path. Middleware are applied in the following order:
//
// 1. Root/global middleware (stored under rootMiddlewarePrefix, e.g. "/") is always added first.
// 2. Any group middleware whose prefix matches the start of the given path.
//
// The returned slice can then be prepended to route-specific handlers when constructing
// the full handler chain for a route. This ensures root middleware always runs first,
// followed by parent/child group middleware.
func (a *App) collectMiddlewares(path string) []Handler {
	var result []Handler

	// Ensure root middleware is added first
	if rootMw, ok := a.middlewares[rootMiddlewarePrefix]; ok {
		result = append(result, rootMw...)
	}

	for prefix, mws := range a.middlewares {
		if prefix == rootMiddlewarePrefix {
			continue // already added root
		}
		if strings.HasPrefix(path, prefix) {
			result = append(result, mws...)
		}
	}

	return result
}
