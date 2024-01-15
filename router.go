package gee

import "strings"

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern parses the pattern, which starts with '/'
func parsePattern(pattern string) []string {
	// 以 / 分割 pattern
	vs := strings.Split(pattern, "/")

	// 用于存储非空的 part
	parts := make([]string, 0)
	for _, item := range vs {
		// 如果 part 非空，则插入 parts
		if item != "" {
			parts = append(parts, item)

			// 如果 part 的第一个字符为 *，则终止循环
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute adds a route to the router
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 以 method 为 key，获取路由树
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 如果路由树不存在，则新建一个
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// 将路由插入到路由树中
	r.roots[method].insert(pattern, parts, 0)

	// 将路由和 handler 存入 handlers 中
	r.handlers[key] = handler
}

// getRoute returns the matching pattern and parameters
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 以 method 为 key，获取路由树
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	// 如果路由树不存在，则返回 nil
	if !ok {
		return nil, nil
	}

	// 从路由树中查找匹配的路由
	n := root.search(searchParts, 0)

	// 如果匹配到路由，则将路由中的参数存入 params
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

// handle finds the matched handler for the given url
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}
