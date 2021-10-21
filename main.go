package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"

	"github.com/becheran/depgraph/pkg/mod"
	"github.com/becheran/depgraph/pkg/server"
	"github.com/becheran/depgraph/pkg/tree"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "localhost:3001", "Set host name or ip address with port were site shall be served")

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Package path argument missing")
		os.Exit(1)
	}
	pkgName := flag.Args()[0]

	js := generateTree(pkgName)

	api := &treeApi{data: js}
	server.Run(host, api.Serve)
}

func generateTree(pkgName string) (jsonOut []byte) {
	cwd, _ := os.Getwd()

	modContent, err := os.ReadFile("go.mod")
	if err != nil {
		log.Fatalf("Unable to read go.mod file. Err: %s", err.Error())
	}
	modInfo := mod.NewModInfo(modContent)

	ctx := &build.Default
	ctx.CgoEnabled = false

	treeBuilder := tree.NewTreeBuilder(ctx, modInfo)
	pkg, err := treeBuilder.DependencyTree(cwd, pkgName)
	if err != nil {
		log.Fatalf("Failed to create dependency tree. Err: %s", err.Error())
	}

	jsonOut, err = json.Marshal(&pkg)
	if err != nil {
		log.Fatalf("Failed to marshal tree to json. Err: %s", err.Error())
	}

	return
}

type treeApi struct {
	data []byte
}

func (api *treeApi) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(api.data)
}
