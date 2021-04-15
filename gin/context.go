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
	StatusCode int

	//处理函数集和下标
	handlers []HandlerFunc
	index    int

	//用来终止后续函数调用
	isOk bool

	//用来启动HtmlTemplate服务
	engine *Engine
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

//调用失败
func (c *Context) Fail(code int, message string) {
	c.JSON(code, H{"message": message})
}

//获取URL中的部分元素
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//返回POST表单中特定的Key对应的value
func (c *Context) PostForm(key string) string {
	return c.R.FormValue(key)
}

//返回URL中key对应的value
func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

//设置状态码并在返回头中设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

///设置状态码
func (c *Context) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

//以string形式返回请求
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, value...)))
}

//以json形式返回请求
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.W)

	if err := encoder.Encode(obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

//以data形式返回请求
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.W.Write(data)
}

//以html返回请求
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)

	if err := c.engine.htmlTemplates.ExecuteTemplate(c.W, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
