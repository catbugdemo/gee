package gee

import "strings"

// Trie

// node 代表每个节点
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// matchChild 用于匹配子节点，第一个匹配成功的节点用于解析参数
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 用于匹配所有子节点，返回所有匹配成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 用于插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归终止条件
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取当前层级的 part
	part := parts[height]
	child := n.matchChild(part)

	// 如果没有匹配到，则新建一个节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// 递归插入下一层级
	child.insert(pattern, parts, height+1)
}

// search 用于搜索节点，找到最长匹配的节点
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 获取当前层级的 part
	part := parts[height]
	children := n.matchChildren(part)

	// 递归搜索下一层级
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
