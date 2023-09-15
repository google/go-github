// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// ProjRootDir returns the go-github root directory that contains dir.
// Returns an error if dir is not in a go-github root.
func ProjRootDir(dir string) (string, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	ok, err := isGoGithubRoot(dir)
	if err != nil {
		return "", err
	}
	if ok {
		return dir, nil
	}
	parent := filepath.Dir(dir)
	if parent == dir {
		return "", fmt.Errorf("not in a go-github root")
	}
	return ProjRootDir(parent)
}

// isGoGithubRoot determines whether dir is the repo root of go-github. It does
// this by checking whether go.mod exists and contains a module directive with
// the path "github.com/google/go-github/vNN".
func isGoGithubRoot(dir string) (bool, error) {
	filename := filepath.Join(dir, "go.mod")
	b, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	mf, err := modfile.ParseLax(filename, b, nil)
	if err != nil {
		// an invalid go.mod file is not a go-github root, so we don't care about the error
		return false, nil
	}
	if mf.Module == nil {
		return false, nil
	}
	// This gets rid of the /vN suffix if it exists.
	base, _, ok := module.SplitPathVersion(mf.Module.Mod.Path)
	if !ok {
		return false, nil
	}
	return base == "github.com/google/go-github", nil
}

type serviceMethod struct {
	receiverName string
	methodName   string
	filename     string
}

func (m *serviceMethod) name() string {
	return fmt.Sprintf("%s.%s", m.receiverName, m.methodName)
}

func getServiceMethods(dir string) ([]*serviceMethod, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var serviceMethods []*serviceMethod
	for _, filename := range dirEntries {
		m, err := getServiceMethodsFromFile(filepath.Join(dir, filename.Name()))
		if err != nil {
			return nil, err
		}
		serviceMethods = append(serviceMethods, m...)
	}
	sort.Slice(serviceMethods, func(i, j int) bool {
		if serviceMethods[i].filename != serviceMethods[j].filename {
			return serviceMethods[i].filename < serviceMethods[j].filename
		}
		return serviceMethods[i].name() < serviceMethods[j].name()
	})
	return serviceMethods, nil
}

// getServiceMethodsFromFile is like getServiceMethodsFromFileDST, but uses
// the AST package instead of the DST package.
func getServiceMethodsFromFile(filename string) ([]*serviceMethod, error) {
	if !strings.HasSuffix(filename, ".go") ||
		strings.HasSuffix(filename, "_test.go") {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Only look at the github package
	if f.Name.Name != "github" {
		return nil, nil
	}
	var serviceMethods []*serviceMethod
	ast.Inspect(f, func(n ast.Node) bool {
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if decl.Recv == nil || len(decl.Recv.List) != 1 {
			return true
		}
		se, ok := decl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		id, ok := se.X.(*ast.Ident)
		if !ok {
			return true
		}
		receiverName := id.Name
		methodName := decl.Name.Name

		// We only want exported methods on exported types.
		// The receiver must either end with Service or be named Client.
		// The exception is github.go, which contains Client methods we want to skip.

		if !ast.IsExported(methodName) || !ast.IsExported(receiverName) {
			return true
		}
		if receiverName != "Client" && !strings.HasSuffix(receiverName, "Service") {
			return true
		}
		if receiverName == "Client" && filepath.Base(filename) == "github.go" {
			return true
		}
		serviceMethods = append(serviceMethods, &serviceMethod{
			receiverName: receiverName,
			methodName:   methodName,
			filename:     filename,
		})
		return true
	})
	return serviceMethods, nil
}