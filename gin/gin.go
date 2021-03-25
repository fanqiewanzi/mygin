package gin

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}


func New() *Engine {
	return &Engine{NewRouter()}
}


func (engine *Engine)GET(pattern string,handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine)POST(pattern string,handler HandlerFunc) {
	engine.router.addRoute("POST", pattern, handler)
}

func (engine *Engine)DELETE(pattern string,handler HandlerFunc) {
	engine.router.addRoute("DELETE", pattern, handler)
}

func (engine *Engine)ServeHTTP(w http.ResponseWriter, r *http.Request){
	c:=NewContext(w,r)
	engine.router.handle(c)
}

func (engine *Engine)Run(addr string) error{
	return http.ListenAndServe(addr,engine)
}