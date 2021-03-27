package gin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//用来接受请求和返回结果的两个对象
	W http.ResponseWriter
	R *http.Request

	//Restful访问方法
	Method string

	//访问完整路径
	Pattern string

	//路径中包含的param元素
	Params map[string]string

	//状态码
	StatueCode int

	//处理函数集和下标
	handlers []HandlerFunc
	index    int

	//用来终止后续函数调用
	isOk bool
}

//Context的构造函数
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:       w,
		R:       r,
		Method:  r.Method,
		Pattern: r.URL.Path,
		index:   -1,
		isOk:    true,
	}
}

//依次执行处理函数集
func (c *Context) Next() {
	//因为会在不同的地方调用Next()函数，保证不能执行到前面执行过的函数,调用range来执行的话会导致死循环
	c.index++
	for length := len(c.handlers); c.index < length; c.index++ {
		if c.isOk {
			c.handlers[c.index](c)
		}
	}
}

//调用Abort终止调用后面的处理函数
func (c *Context) Abort() {
	c.isOk = false
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.R.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) Statue(code int) {
	c.StatueCode = code
	c.W.WriteHeader(code)
}

func (c *Context) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Statue(code)
	c.W.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Statue(code)
	encoder := json.NewEncoder(c.W)

	if err := encoder.Encode(obj); err != nil {
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
