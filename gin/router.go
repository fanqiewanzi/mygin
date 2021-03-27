package gin

import (
	"net/http"
	"strings"
)

type router struct {

	//动态路由树
	roots map[string]*node

	//路由map
	handlers map[string]HandlerFunc
}

//构造路由
func NewRouter() *router {
	return &router{make(map[string]*node), make(map[string]HandlerFunc)}
}

//解析路由路径
func parsePattern(pattern string) []string {

	//将路径由'/'分开，返回一个切片，里面存储的是路径的各个元素
	str := strings.Split(pattern, "/")

	//创建一个切片用来存储路径的各部分
	paths := make([]string, 0)

	//当扫描到的路径元素不为空时进行复制，扫描到文件开头的字符时退出扫描，因为文件后面的路径都是无效路径
	for _, path := range str {
		if path != "" {
			paths = append(paths, path)
			if path[0] == '*' {
				break
			}
		}
	}

	return paths
}

//增加路由路径
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {

	//分析具体路径
	paths := parsePattern(pattern)

	//将访问方法和路径结合
	key := method + "-" + pattern

	//在动态路由树中查找是否存在这条路径下对应的方法(GET/POST/DELETE...)
	//如果不存在则在对应的路由树下新建一个空的节点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	//路径初始化
	r.roots[method].insert(pattern, paths, 0)

	//将路径和处理函数建立关系
	r.handlers[key] = handler
}

//查找路由中的元素
func (r *router) getRoute(method string, path string) (*node, map[string]string) {

	//解析路由将其分为一个切片
	paths := parsePattern(path)

	//创建params切片用来存储URL中的param元素
	params := make(map[string]string)

	//检查访问方法(GET/DELETE/POST/PUT)是否存在，不存在则不存在此路由
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	//检查是否存在此路由
	n := root.search(paths, 0)

	//路由存在
	if n != nil {
		//解析全路径成为一个切片
		parts := parsePattern(n.pattern)
		//将切片中的param和动态路由信息解析出来并保存
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = paths[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(paths[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {

	//查找路由解析路径
	n, params := r.getRoute(c.Method, c.Pattern)

	//运行相应处理函数
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Pattern)
		})
	}
	c.Next()
}
