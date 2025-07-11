package nodes

import (
	"encoding/json"
	"reflect"

	"github.com/tylertravisty/go-lexical"
)

var _ lexical.Node = &ParagraphNode{}

// ParagraphNode implements the lexical paragraph node type
type ParagraphNode struct {
	ElementNode
	TextFormat int    `json:"textFormat"`
	TextStyle  string `json:"textStyle"`
}

// Find saves paragraph node to nodes if paragraph type is in map and then calls find on children
func (pn *ParagraphNode) Find(nodes map[string][]lexical.Node) {
	Find(pn, nodes)

	for _, child := range pn.Children {
		child.Find(nodes)
	}
}

// Type returns type of paragraph node
func (pn ParagraphNode) Type() (string, reflect.Type) {
	return "paragraph", reflect.TypeOf(pn)
}

// Unmarshal unmarshals the paragraph node
func (pn *ParagraphNode) Unmarshal(data map[string]interface{}) error {
	pnB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(pnB, pn)
}
