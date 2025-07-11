package nodes

import "github.com/tylertravisty/go-lexical"

const (
	pkg = "nodes"
)

// Find is a helper function to save the node in nodes if the node type exists in the map
func Find(node lexical.Node, nodes map[string][]lexical.Node) {
	nodeType, _ := node.Type()
	if save, exists := nodes[nodeType]; exists {
		save = append(save, node)
		nodes[nodeType] = save
	}
}
