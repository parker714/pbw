package pbw

import (
	"reflect"
	"strings"
	"testing"
)

func TestTrie_insert(t *testing.T) {
	n := newNode("")
	n.insert([]string{"/", "user"})
	n.insert([]string{"/", "sku"})

	except := &node{
		isWild: strings.HasPrefix("", strColon) || strings.HasPrefix("", strStar),
		next: map[string]*node{
			"/": {
				next: map[string]*node{
					"user": {
						next:     make(map[string]*node),
						patterns: []string{"/", "user"},
					},
					"sku": {
						next:     make(map[string]*node),
						patterns: []string{"/", "sku"},
					},
				},
				patterns: make([]string, 0),
			},
		},
		patterns: make([]string, 0),
	}

	if !reflect.DeepEqual(except, n) {
		t.Fatalf("except %v, got %v", except, n)
	}
}

func TestTrie_search(t *testing.T) {
	n := newNode("")
	n.insert([]string{"/", "user"})
	n.insert([]string{"/", ":lang", "doc"})

	except := []string{"/", ":lang", "doc"}
	got := n.search([]string{"/", "go", "doc"}, 0)
	if !reflect.DeepEqual(except, got) {
		t.Fatalf("except %v, got %v", except, got)
	}

	except = []string{}
	got = n.search([]string{"/", "go", "doc1"}, 0)
	if !reflect.DeepEqual(except, got) {
		t.Fatalf("except %v, got %v", except, got)
	}
}
