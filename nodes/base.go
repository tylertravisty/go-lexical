package nodes

import "encoding/json"

// BaseNode is the basic node type
type BaseNode struct {
	NodeType string `json:"type"`
	Version  int    `json:"version"`
}

// Unmarshal unmarshals a base node
func (bn *BaseNode) Unmarshal(data map[string]any) error {
	bnB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(bnB, bn)
}
