package gin

import "strings"

// FIXME:这个结构体里边的字段为啥不换行了，保持风格一致。
type node struct {
	//路由部分路径
	path string
	//路由全路径
	pattern string
	//子节点
	children []*node
	//是否精确匹配
	isWild bool
}

//单个子路径元素是否存在，存在返回元素节点
func (n *node) matchChild(path string) *node {

	//从根节点开始找，当子元素路径相同或找到动态子元素时返回元素节点
	for _, child := range n.children {
		if child.path == path || child.isWild {
			return child
		}
	}
	//找不到返回nil
	return nil
}

//构造路径，返回一个Node切片
func (n *node) matchChildren(path string) []*node {

	//创建一个新的Node切片
	nodes := make([]*node, 0)

	//从根节点的孩子节点中查找部分路径，找到就放入Node切片中
	for _, child := range n.children {
		if child.path == path || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//路由树增加路由
func (n *node) insert(pattern string, paths []string, height int) {

	//如果找的长度等于路径的高度，代表已经找到路径尽头，此时节点的全路径可能发生变化，将当前节点的全路径赋值并返回
	if len(paths) == height {
		n.pattern = pattern
		return
	}

	//从全路径的最开始的元素找是否存在该节点
	path := paths[height]
	child := n.matchChild(path)
	//不存在则创建新节点
	if child == nil {
		child = &node{path: path, isWild: path[0] == ':' || path[0] == '*'}
		n.children = append(n.children, child)
	}
	//递归直到走完全路径长度，就创建了一条新的路由
	child.insert(pattern, paths, height+1)
}

//搜索路由
func (n *node) search(paths []string, height int) *node {

	//找到路径最大高度或者找到通配符，当完整路径不为空字符串就返回最后的子节点
	if len(paths) == height || strings.HasPrefix(n.path, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	path := paths[height]
	children := n.matchChildren(path)

	//从分支节点中找到对应的路径，找到路径就返回
	for _, child := range children {
		result := child.search(paths, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
