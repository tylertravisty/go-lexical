package nodes

import (
	"encoding/json"
	"reflect"

	"github.com/tylertravisty/go-lexical"
)

var _ lexical.Node = &AutoLinkNode{}

// AutoLinkNode implements the lexical autolink node type
type AutoLinkNode struct {
	LinkNode
	IsUnlinked bool `json:"isUnlinked"`
}

// Find saves autolink node to nodes if autolink type is in map and then calls find on children
func (aln *AutoLinkNode) Find(nodes map[string][]lexical.Node) {
	Find(aln, nodes)

	for _, child := range aln.Children {
		child.Find(nodes)
	}
}

// Type returns type of autolink node
func (aln AutoLinkNode) Type() (string, reflect.Type) {
	return "autolink", reflect.TypeOf(aln)
}

// Unmarshal unmarshals the autolink node
func (aln *AutoLinkNode) Unmarshal(data map[string]interface{}) error {
	alnB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(alnB, aln)
}
