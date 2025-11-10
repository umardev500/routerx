package routerx

type Handler func(*Ctx) error

type Router interface {
	Delete(path string, handlers ...Handler)
	Get(path string, handlers ...Handler)
	Group(prefix string) Router
	Patch(path string, handlers ...Handler)
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Use(handlers ...Handler)
}
