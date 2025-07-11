package nodes

import "github.com/tylertravisty/go-lexical"

// RootNode is a lexical root node
type RootNode struct {
	Root ElementNode `json:"root"`
}

// Find finds the specified node types
func (rn *RootNode) Find(nodes map[string][]lexical.Node) {
	rn.Root.Find(nodes)
}
