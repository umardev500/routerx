package routerx

import "net/http"

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
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		c := NewCtx(w, r, handlers)

		if r.Method != method {
			c.Status(http.StatusMethodNotAllowed).
				Json(http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		if err := c.Next(); err != nil {
			c.Status(http.StatusInternalServerError).
				Json(err.Error())
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
	NewGroup(prefix, a)
	return a
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
