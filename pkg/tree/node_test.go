package tree_test

import (
	"encoding/json"
	"testing"

	"github.com/becheran/depgraph/pkg/tree"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJson(t *testing.T) {
	var suite = []struct {
		in  *tree.PackageNode
		out string
	}{
		{&tree.PackageNode{ID: "Foo"}, `[{"data":{"id":"Foo"}}]`},
		{&tree.PackageNode{ID: ""}, `[]`},
		{&tree.PackageNode{ID: "Foo", Imports: []*tree.PackageNode{
			{ID: "Bar"},
		}}, `[{"data":{"id":"Foo"}},{"data":{"id":"Bar"}},{"data":{"source":"Foo","target":"Bar"}}]`},
	}
	for _, s := range suite {
		t.Run(s.out, func(t *testing.T) {
			res, err := json.Marshal(s.in)
			assert.Nil(t, err)
			assert.Equal(t, s.out, string(res))
		})
	}
}

func TestTraverse(t *testing.T) {
	const expectedName = "First"

	nodes := tree.PackageNode{
		ID: expectedName,
	}

	called := false
	nodes.Traverse(func(node *tree.PackageNode) {
		if called {
			t.Fatal("Expected callback to be called only once")
		}
		called = true
		assert.Equal(t, expectedName, node.ID)
	})
}

func TestTraverseComplex(t *testing.T) {
	expectedNames := []string{"root", "foo", "bar", "baz"}

	root := tree.PackageNode{
		ID: expectedNames[0],
		Imports: []*tree.PackageNode{
			{ID: expectedNames[1]},
			{ID: expectedNames[2], Imports: []*tree.PackageNode{
				{ID: expectedNames[3]},
			}},
		},
	}

	ctr := 0
	root.Traverse(func(node *tree.PackageNode) {
		if ctr >= len(expectedNames) {
			t.Fatal("Callback called to often")
		}
		assert.Equal(t, expectedNames[ctr], node.ID)
		ctr++
	})

	assert.Equal(t, len(expectedNames), ctr)
}
