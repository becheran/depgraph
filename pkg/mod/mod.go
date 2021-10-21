package mod

import (
	"log"
	"strings"

	"golang.org/x/mod/modfile"
)

type ModInfos interface {
	Type(path string) ImportType
}

type ImportType string

const (
	ImportTypeUnknown = ""
	ImportTypeStdLib  = "STD"
	ImportTypePkg     = "PKG"
	ImportTypeReq     = "REQ"
	ImportTypeCgo     = "CGO"
)

type ModInfo struct {
	pkgPrefix        string
	requiredPrefixes []string
}

func NewModInfo(modContent []byte) *ModInfo {
	file, err := modfile.Parse("go.mod", modContent, nil)
	if err != nil {
		log.Fatalf("Failed to parse go.mod file. Err: %s", err.Error())
	}
	req := make([]string, len(file.Require))
	for idx, r := range file.Require {
		req[idx] = r.Mod.Path
	}
	return &ModInfo{
		pkgPrefix:        file.Module.Mod.Path,
		requiredPrefixes: req,
	}
}

func (inf *ModInfo) Type(path string) ImportType {
	if strings.HasPrefix(path, inf.pkgPrefix) {
		return ImportTypePkg
	}
	for _, req := range inf.requiredPrefixes {
		if strings.HasPrefix(path, req) {
			return ImportTypeReq
		}
	}
	// TODO: check if is cgo pkg
	return ImportTypeStdLib
}
