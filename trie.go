package pbw

import "strings"

var (
	strColon = ":"
	strStar  = "*"
)

type node struct {
	isWild bool
	next   map[string]*node

	patterns []string
}

func newNode(part string) *node {
	return &node{
		isWild:   strings.HasPrefix(part, strColon) || strings.HasPrefix(part, strStar),
		next:     make(map[string]*node),
		patterns: make([]string, 0),
	}
}

func (n *node) insert(patterns []string) {
	curr := n
	for _, part := range patterns {
		if _, ok := curr.next[part]; !ok {
			curr.next[part] = newNode(part)
		}
		curr = curr.next[part]
	}
	curr.patterns = patterns
}

func (n *node) search(reqParts []string, index int) []string {
	if len(reqParts) == index {
		return n.patterns
	}

	for part, next := range n.next {
		if next.isWild || reqParts[index] == part {
			return next.search(reqParts, index+1)
		}
	}
	return make([]string, 0)
}
