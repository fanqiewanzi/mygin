package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string] interface{}

type Context struct {
	W http.ResponseWriter
	R *http.Request

	Method string
	Path string
	Params map[string]string

	StatueCode int
}

func NewContext(w http.ResponseWriter,r *http.Request) *Context  {
	return &Context{W:w,
		R:r,
		Method:r.Method,
		Path:r.URL.Path}
}

func (c *Context)Param(key string) string {
	value,_:=c.Params[key]
	return value
}

func (c *Context)PostForm(key string) string  {
	return c.R.FormValue(key)
}

func (c *Context)Query(key string) string  {
	return c.R.URL.Query().Get(key)
}

func (c *Context) Statue (code int)  {
	c.StatueCode=code
	c.W.WriteHeader(code)
}

func (c *Context) SetHeader (key,value string)  {
	c.W.Header().Set(key,value)
}

func (c *Context) String(code int,format string,value...interface{})  {
	c.SetHeader("Content-Type","text/plain")
	c.Statue(code)
	c.W.Write([]byte(fmt.Sprintf(format,value...)))
}

func (c *Context)JSON(code int,obj interface{})  {
	c.SetHeader("Content-Type","application/json")
	c.Statue(code)
	encoder:=json.NewEncoder(c.W)

	if err:=encoder.Encode(obj);err!=nil{
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Statue(code)
	c.W.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Statue(code)
	c.W.Write([]byte(html))
}
