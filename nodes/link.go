package nodes

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tylertravisty/go-lexical"
)

var _ lexical.Node = &LinkNode{}

// LinkNode implements the lexical link node type
type LinkNode struct {
	ElementNode
	Rel    *string `json:"rel"`
	Target *string `json:"target"`
	Title  *string `json:"title"`
	URL    string  `json:"url"`
}

// Find saves link node to nodes if link type is in map and then calls find on children
func (ln *LinkNode) Find(nodes map[string][]lexical.Node) {
	Find(ln, nodes)

	for _, child := range ln.Children {
		child.Find(nodes)
	}
}

// Type returns type of link node
func (ln LinkNode) Type() (string, reflect.Type) {
	return "link", reflect.TypeOf(ln)
}

// Unmarshal unmarshals the link node
func (ln *LinkNode) Unmarshal(data map[string]interface{}) error {
	lnB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(lnB, ln)
}

// Valid verifies the link node is valid
func (ln *LinkNode) Valid() error {
	err := ln.ElementNode.Valid()
	if err != nil {
		return err
	}

	err = runLinkNodeValFuncs(
		ln,
		linkNodeRequireTextChild,
	)
	if err != nil {
		return fmt.Errorf("%s: invalid link node: %v", pkg, err)
	}

	return nil
}

type linkNodeValFunc func(*LinkNode) error

func runLinkNodeValFuncs(node *LinkNode, fns ...linkNodeValFunc) error {
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

func linkNodeRequireTextChild(node *LinkNode) error {
	if len(node.Children) != 1 {
		return fmt.Errorf("invalid number of children")
	}

	child := node.Children[0]
	cType, _ := child.Type()
	expectedType, _ := TextNode{}.Type()
	if cType != expectedType {
		return fmt.Errorf("invalid child type")
	}

	return nil
}
