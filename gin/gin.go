package gin

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {

	//继承GroupRouter，拥有其所有方法
	*GroupRouter

	//Engine的路由树
	router *router

	//Engine拥有的所有组集
	groups []*GroupRouter
}

type GroupRouter struct {

	//组的前缀名
	prefix string

	//组使用中间件的处理函数集
	middleWares []HandlerFunc

	//每个组都使用同一个Engine实例
	engine *Engine
}

//Engine的构造函数
func New() *Engine {
	e := &Engine{router: NewRouter()}
	e.GroupRouter = &GroupRouter{engine: e}
	e.groups = append(e.groups, e.GroupRouter)
	return e
}

//组的构造函数，将前缀和engine指针设好并将组指针加入到engine中的组集中
func (g *GroupRouter) Group(prefix string) *GroupRouter {
	e := g.engine
	newGroup := &GroupRouter{prefix: prefix, engine: e}
	e.groups = append(e.groups, newGroup)
	return newGroup
}

func (g *GroupRouter) Use(handlerFunc ...HandlerFunc) {
	g.middleWares = append(g.middleWares, handlerFunc...)
}

//在原有的路由下加入组前缀再添加路由
func (g *GroupRouter) addRoute(method string, path string, handler HandlerFunc) {
	pattern := g.prefix + path
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *GroupRouter) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *GroupRouter) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

func (g *GroupRouter) DELETE(pattern string, handler HandlerFunc) {
	g.addRoute("DELETE", pattern, handler)
}

func (g *GroupRouter) PUT(pattern string, handler HandlerFunc) {
	g.addRoute("PUT", pattern, handler)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var middleWares []HandlerFunc

	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}

	c := NewContext(w, r)
	c.handlers = middleWares
	engine.router.handle(c)
}

func (engine *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
