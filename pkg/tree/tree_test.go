package tree_test

import (
	"fmt"
	"go/build"
	"sync/atomic"
	"testing"

	"github.com/becheran/depgraph/pkg/mod"
	"github.com/becheran/depgraph/pkg/tree"
	"github.com/stretchr/testify/assert"
)

type ctx struct {
	OnImport func(path string, srcDir string, mode build.ImportMode) (*build.Package, error)
}

func NewTestCtx() *ctx {
	return &ctx{
		OnImport: func(path, srcDir string, mode build.ImportMode) (*build.Package, error) { return &build.Package{}, nil },
	}
}

func (ctx *ctx) Import(path string, srcDir string, mode build.ImportMode) (*build.Package, error) {
	return ctx.OnImport(path, srcDir, mode)
}

type importType struct {
	retVal mod.ImportType
}

func NewImportType() *importType {
	return &importType{
		retVal: mod.ImportTypePkg,
	}
}

func (it *importType) Type(s string) mod.ImportType {
	return it.retVal
}

func TestNoPackage(t *testing.T) {
	ctx := NewTestCtx()

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())
	_, err := treeBuilder.DependencyTree("", "")

	assert.Nil(t, err)
}

func TestFailingCtx_err(t *testing.T) {
	ctx := NewTestCtx()
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		return nil, fmt.Errorf("Unknown")
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())
	_, err := treeBuilder.DependencyTree("", "")

	assert.NotNil(t, err)
}

func TestSimpleTree(t *testing.T) {
	ctx := NewTestCtx()
	packages := []build.Package{
		{ImportPath: "P0", Imports: []string{"P1"}},
		{ImportPath: "P1", Imports: []string{"P2"}},
		{ImportPath: "P2", Imports: []string{"P3"}},
		{ImportPath: "P3"},
	}

	idx := 0
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		pkg := packages[idx]
		idx++
		return &pkg, nil
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())
	root, err := treeBuilder.DependencyTree("", "")
	assert.Nil(t, err)

	travIdx := 0
	root.Traverse(func(node *tree.PackageNode) {
		assert.Equal(t, packages[travIdx].ImportPath, node.ID)
		travIdx++
	})

	assert.Equal(t, travIdx, len(packages))
}

func TestRecurseReturnErr(t *testing.T) {
	ctx := NewTestCtx()
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		return &build.Package{Name: "SameNameAgain"}, nil
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())
	_, err := treeBuilder.DependencyTree("", "")

	assert.NotNil(t, err)
}

func TestRecurseExceedMaxDepthErr(t *testing.T) {
	ctx := NewTestCtx()
	idx := 0
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		p := build.Package{ImportPath: fmt.Sprintf("P%d", idx), Imports: []string{fmt.Sprintf("P%d", idx+1)}}
		idx++
		if idx > tree.MaxDepth {
			assert.FailNow(t, "Exceeded max depth for tree iteration")
		}
		return &p, nil
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())
	_, err := treeBuilder.DependencyTree("", "")

	assert.NotNil(t, err)
}

func BenchmarkDepthTree(b *testing.B) {
	ctx := NewTestCtx()
	idx := 0
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		p := build.Package{ImportPath: fmt.Sprintf("P%d", idx), Imports: []string{fmt.Sprintf("P%d", idx+1)}}
		idx++
		if idx > 200 {
			p.Imports = nil
		}
		return &p, nil
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())

	for n := 0; n < b.N; n++ {
		treeBuilder.DependencyTree("", "")
	}
}

func BenchmarkBreathTree(b *testing.B) {
	ctx := NewTestCtx()
	var idx int64 = 0
	ctx.OnImport = func(path, srcDir string, mode build.ImportMode) (*build.Package, error) {
		newIdx := atomic.AddInt64(&idx, 1)
		p := build.Package{
			ImportPath: path, // TODO Deeper
			Imports:    []string{fmt.Sprintf("A%d", newIdx+1), fmt.Sprintf("B%d", newIdx+1), fmt.Sprintf("C%d", newIdx+1)},
		}
		if newIdx > 200 {
			p.Imports = nil
		}
		return &p, nil
	}

	treeBuilder := tree.NewTreeBuilder(ctx, NewImportType())

	for n := 0; n < b.N; n++ {
		treeBuilder.DependencyTree("", "")
	}
}
