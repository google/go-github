package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
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

func ExitErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

type ServiceMethod struct {
	ReceiverName string
	MethodName   string
	Filename     string
}

func (m *ServiceMethod) Name() string {
	return fmt.Sprintf("%s.%s", m.ReceiverName, m.MethodName)
}

func GetServiceMethods(dir string) ([]*ServiceMethod, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var serviceMethods []*ServiceMethod
	for _, filename := range dirEntries {
		m, err := getServiceMethodsFromFile(filepath.Join(dir, filename.Name()))
		if err != nil {
			return nil, err
		}
		serviceMethods = append(serviceMethods, m...)
	}
	sort.Slice(serviceMethods, func(i, j int) bool {
		if serviceMethods[i].Filename != serviceMethods[j].Filename {
			return serviceMethods[i].Filename < serviceMethods[j].Filename
		}
		return serviceMethods[i].Name() < serviceMethods[j].Name()
	})
	return serviceMethods, nil
}

func getServiceMethodsFromFile(filename string) ([]*ServiceMethod, error) {
	if !strings.HasSuffix(filename, ".go") ||
		strings.HasSuffix(filename, "_test.go") {
		return nil, nil
	}

	df, err := decorator.ParseFile(nil, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	// Only look at the github package
	if df.Name.Name != "github" {
		return nil, nil
	}
	var serviceMethods []*ServiceMethod
	dst.Inspect(df, func(n dst.Node) bool {
		decl, ok := n.(*dst.FuncDecl)
		if !ok {
			return true
		}
		if decl.Recv == nil || len(decl.Recv.List) != 1 {
			return true
		}
		se, ok := decl.Recv.List[0].Type.(*dst.StarExpr)
		if !ok {
			return true
		}
		id, ok := se.X.(*dst.Ident)
		if !ok {
			return true
		}
		receiverName := id.Name
		methodName := decl.Name.Name
		if !dst.IsExported(methodName) || !dst.IsExported(receiverName) {
			return true
		}
		if receiverName != "Client" && !strings.HasSuffix(receiverName, "Service") {
			return true
		}
		// Skip Client methods in github.go
		if receiverName == "Client" && filepath.Base(filename) == "github.go" {
			return true
		}
		serviceMethods = append(serviceMethods, &ServiceMethod{
			ReceiverName: receiverName,
			MethodName:   methodName,
			Filename:     filename,
		})
		return true
	})
	return serviceMethods, nil
}
