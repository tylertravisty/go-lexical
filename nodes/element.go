package nodes

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tylertravisty/go-lexical"
)

var _ lexical.Node = &ElementNode{}

// ElementNode implements the lexical element node type
type ElementNode struct {
	BaseNode
	Children  lexical.NodeArray `json:"children"`
	Direction *string           `json:"direction"`
	Format    string            `json:"format"`
	Indent    int               `json:"indent"`
}

// Find saves element node to nodes if element type is in map and then calls find on children
func (en *ElementNode) Find(nodes map[string][]lexical.Node) {
	Find(en, nodes)

	for _, child := range en.Children {
		child.Find(nodes)
	}
}

// TextContentSize returns the text content size of the element's children
func (en ElementNode) TextContentSize() int {
	size := 0
	for _, node := range en.Children {
		size = size + node.TextContentSize()
	}

	return size
}

// Type returns type of element node
func (en ElementNode) Type() (string, reflect.Type) {
	return "element", reflect.TypeOf(en)
}

// Unmarshal unmarshals the element node
func (en *ElementNode) Unmarshal(data map[string]any) error {
	enB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(enB, en)
}

// Valid verifies the element node is valid
func (en *ElementNode) Valid() error {
	err := runElementNodeValFuncs(
		en,
		elementNodeRequireDirection,
		elementNodeRequireFormat,
		elementNodeChildrenValid,
	)
	if err != nil {
		return fmt.Errorf("%s: invalid element node: %v", pkg, err)
	}

	return nil
}

type elementNodeValFunc func(*ElementNode) error

func runElementNodeValFuncs(node *ElementNode, fns ...elementNodeValFunc) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	for _, fn := range fns {
		err := fn(node)
		if err != nil {
			return err
		}
	}

	return nil
}

func elementNodeRequireDirection(node *ElementNode) error {
	if node.Direction == nil {
		return nil
	}

	switch *node.Direction {
	case "ltr", "rtl":
	default:
		return fmt.Errorf("invalid direction")
	}

	return nil
}

func elementNodeRequireFormat(node *ElementNode) error {
	switch node.Format {
	case "left", "start", "center", "right", "end", "justify", "":
	default:
		return fmt.Errorf("invalid format")
	}

	return nil
}

func elementNodeChildrenValid(node *ElementNode) error {
	for _, child := range node.Children {
		err := child.Valid()
		if err != nil {
			return err
		}
	}

	return nil
}
