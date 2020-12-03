package pbw

import "testing"

func mockNode() Node {
	node := NewNode()

	parts := [][]string{
		{},
		{"welcome"},
		{"user", ":name"},
		{"user", ":name", "*action"},
	}
	for _, v := range parts {
		node.Insert(v)
	}
	return node
}

func TestNode_Search(t *testing.T) {
	node := mockNode()

	cases := []struct {
		input  string
		expect string
	}{
		{
			"/welcome",
			"welcome",
		},
		{
			"/user/pb",
			"user/:name",
		},
		{
			"/user/pb/say",
			"user/:name/*action",
		},
		{
			"/book/pb",
			"",
		},
	}
	for _, c := range cases {
		if pattern := node.Search(c.input); c.expect != pattern {
			t.Fatalf("expect %s, got %s, input %s", c.expect, pattern, c.input)
		}
	}
}
