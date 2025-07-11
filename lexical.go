package lexical

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

const (
	pkg = "lexical"
)

// TypeMap holds the registered lexical node types
type TypeMap struct {
	mu    sync.RWMutex
	types map[string]reflect.Type
}

// Add associates a lexical node type with the given name
func (tm *TypeMap) Add(name string, nodeType reflect.Type) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.types[name] = nodeType
}

// Type returns the lexical node type associated with given name
func (tm *TypeMap) Type(name string) (reflect.Type, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	t, exists := tm.types[name]
	return t, exists
}

// DefaultNodeTypes is a global type map for lexical nodes
var DefaultNodeTypes = &TypeMap{types: map[string]reflect.Type{}}

// ResetNodes resets the default global node type map
func ResetNodes() {
	DefaultNodeTypes.mu.Lock()
	defer DefaultNodeTypes.mu.Unlock()

	DefaultNodeTypes = &TypeMap{types: map[string]reflect.Type{}}
}

// Node defines the interface for lexical nodes
type Node interface {
	Find(nodes map[string][]Node)
	TextContentSize() int
	Type() (string, reflect.Type)
	Unmarshal(data map[string]interface{}) error
	Valid() error
}

// NodeArray is an array of Nodes
type NodeArray []Node

// RegisterNode registers the lexical node
func RegisterNode(node Node) error {
	name, nodeType := node.Type()
	if _, exists := DefaultNodeTypes.Type(name); exists {
		return fmt.Errorf("%s: node type already exists: %v", pkg, name)
	}

	DefaultNodeTypes.Add(name, nodeType)
	return nil
}

// RegisterNodes registers the lexical nodes
func RegisterNodes(nodes ...Node) error {
	for _, node := range nodes {
		err := RegisterNode(node)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal unmarshals map into node
func Unmarshal(data map[string]any) (Node, error) {
	nodeTypeName, ok := data["type"].(string)
	if !ok {
		return nil, fmt.Errorf("%s: invalid node type", pkg)
	}

	nodeType, exists := DefaultNodeTypes.Type(nodeTypeName)
	if !exists {
		return nil, fmt.Errorf("%s: unsupported node type", pkg)
	}

	node, ok := reflect.New(nodeType).Interface().(Node)
	if !ok {
		return nil, fmt.Errorf("%s: invalid node", pkg)
	}

	err := node.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s: error unmarshaling node: %v", pkg, err)
	}

	return node, nil
}

// UnmarshalJSON unmarshals bytes into node array
func (na *NodeArray) UnmarshalJSON(data []byte) error {
	var obj []map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	var array []Node
	var node Node
	for _, data := range obj {
		node, err = Unmarshal(data)
		if err != nil {
			return err
		}

		array = append(array, node)
	}

	*na = array

	return nil
}
