package pbw

import (
	"strings"
)

type Node interface {
	Insert(parts []string)
	Search(reqPath string) string
}

type node struct {
	pattern string
	isWild  bool
	next    map[string]*node // map[part]*node
}

func NewNode() Node {
	return &node{next: make(map[string]*node)}
}

func (n *node) Insert(parts []string) {
	if len(parts) == 0 {
		return
	}

	curr := n
	for _, part := range parts {
		if _, ok := curr.next[part]; !ok {
			curr.next[part] = &node{
				next:   make(map[string]*node),
				isWild: strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*"),
			}
		}
		curr = curr.next[part]
	}
	curr.pattern = strings.Join(parts, "/")
}

// Search for route in node, return pattern
func (n *node) Search(reqPath string) string {
	return search(n, strings.Split(reqPath, "/")[1:], 0)
}

func search(n *node, reqParts []string, index int) string {
	if len(reqParts) == index {
		return n.pattern
	}

	for part, next := range n.next {
		if reqParts[index] == part || next.isWild {
			return search(next, reqParts, index+1)
		}
	}
	return ""
}
