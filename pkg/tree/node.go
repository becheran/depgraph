package tree

import (
	"encoding/json"
	"sort"

	"github.com/becheran/depgraph/pkg/mod"
)

type PackageNode struct {
	ID      string
	Name    string
	Type    mod.ImportType
	Imports []*PackageNode
}

type JsonNodeData struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Source string `json:"source,omitempty"`
	Target string `json:"target,omitempty"`
}

type JsonNode struct {
	Data JsonNodeData `json:"data,omitempty"`
}

func (node *PackageNode) MarshalJSON() (res []byte, err error) {
	nodeSet := make(map[string]*PackageNode)
	node.Traverse(func(current *PackageNode) {
		if current.ID == "" {
			// Ignore empty nodes
			return
		}
		nodeSet[current.ID] = current
	})

	connections := make([]JsonNode, 0, len(nodeSet))
	nodes := make([]JsonNode, 0, len(nodeSet))
	for nodeID, node := range nodeSet {
		nodes = append(nodes, JsonNode{
			Data: JsonNodeData{
				ID:   nodeID,
				Name: node.Name,
				Type: string(node.Type),
			},
		})
		for _, imp := range node.Imports {
			// TODO optimize alloc with one grow call
			connections = append(connections, JsonNode{
				Data: JsonNodeData{
					Source: nodeID,
					Target: imp.ID,
				},
			})
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Data.ID > nodes[j].Data.ID
	})

	combined := append(nodes, connections...)
	return json.Marshal(&combined)
}

func (node *PackageNode) Traverse(cb func(node *PackageNode)) {
	traverse(node, cb)
	return
}

func traverse(current *PackageNode, cb func(node *PackageNode)) {
	if current == nil {
		return
	}
	cb(current)
	for _, imp := range current.Imports {
		traverse(imp, cb)
	}
	return
}
