package nodes

import (
	"encoding/json"
	"reflect"

	"github.com/tylertravisty/go-lexical"
)

var _ lexical.Node = &TextNode{}

// TextNode implements the lexical text node type
type TextNode struct {
	BaseNode
	Detail int    `json:"detail"`
	Format int    `json:"format"`
	Mode   string `json:"mode"`
	Style  string `json:"style"`
	Text   string `json:"text"`
}

// Find saves text node to nodes if text type is in map
func (tn *TextNode) Find(nodes map[string][]lexical.Node) {
	Find(tn, nodes)
}

// TextContentSize returns the length of the text
func (tn *TextNode) TextContentSize() int {
	return len(tn.Text)
}

// Type returns type of text node
func (tn TextNode) Type() (string, reflect.Type) {
	return "text", reflect.TypeOf(tn)
}

// Unmarshal unmarshals the text node
func (tn *TextNode) Unmarshal(data map[string]interface{}) error {
	tnB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(tnB, tn)
}

// Valid verifies the text node is valid
func (tn *TextNode) Valid() error {
	return nil
}
