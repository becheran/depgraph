package tree

import (
	"fmt"
	"go/build"
	"sync"

	"github.com/becheran/depgraph/pkg/mod"
)

// MaxDepth for tree traversal
const MaxDepth = 100

type BuildCtx interface {
	Import(path string, srcDir string, mode build.ImportMode) (*build.Package, error)
}

type TreeBuilder struct {
	ctx      BuildCtx
	modInfo  mod.ModInfos
	cacheMux sync.Mutex
	cache    map[string]*PackageNode
}

func NewTreeBuilder(ctx BuildCtx, modInfo mod.ModInfos) *TreeBuilder {
	return &TreeBuilder{
		ctx:     ctx,
		modInfo: modInfo,
	}
}

func (tb *TreeBuilder) DependencyTree(dir, pkgName string) (node *PackageNode, err error) {
	// Reset cache
	tb.cache = make(map[string]*PackageNode)

	node, err = tb.recurseTree(dir, pkgName, 0)
	if err != nil {
		return
	}

	return
}

func (tb *TreeBuilder) recurseTree(root, pkgName string, depth int) (node *PackageNode, err error) {
	tb.cacheMux.Lock()
	if cacheNode, exists := tb.cache[pkgName]; exists {
		tb.cacheMux.Unlock()
		node = cacheNode
		return
	}
	tb.cache[pkgName] = nil // Reserve empty node
	tb.cacheMux.Unlock()

	removeNode := func() {
		tb.cacheMux.Lock()
		delete(tb.cache, pkgName)
		tb.cacheMux.Unlock()
	}

	pkg, err := tb.ctx.Import(pkgName, root, 0)
	if err != nil {
		removeNode()
		return
	}

	// TODO: Configure if packages in std lib shall be visited
	if pkg.Goroot {
		removeNode()
		return
	}

	depth++
	if depth >= MaxDepth {
		removeNode()
		err = fmt.Errorf("max tree depth of %d exceeded", MaxDepth)
		return
	}

	tb.cacheMux.Lock()
	node = &PackageNode{
		ID:   pkg.ImportPath,
		Name: pkg.Name,
		Type: tb.modInfo.Type(pkg.ImportPath),
	}
	tb.cache[pkgName] = node
	tb.cacheMux.Unlock()

	if node.Type == mod.ImportTypeReq {
		// Do not resolve import dependencies
		return
	}

	var importMux sync.Mutex
	node.Imports = make([]*PackageNode, 0, len(pkg.Imports))
	var importErr error
	var wg sync.WaitGroup
	for _, imported := range pkg.Imports {
		wg.Add(1)
		importedCpy := imported
		go func() {
			defer wg.Done()
			next, err := tb.recurseTree(pkg.Dir, importedCpy, depth)
			importMux.Lock()
			if err != nil {
				importErr = err
			}
			if next != nil {
				node.Imports = append(node.Imports, next)
			}
			importMux.Unlock()
		}()
	}
	wg.Wait()
	err = importErr

	return
}
