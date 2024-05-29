package scanner

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// ScanDir scans the given directory for Go files and returns a list of file paths.
func ScanDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ParseFile parses the Go file and returns a list of handler functions.
func ParseFile(filename string) ([]*ast.FuncDecl, *token.FileSet, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var handlers []*ast.FuncDecl
	for _, f := range node.Decls {
		if fn, isFn := f.(*ast.FuncDecl); isFn {
			if fn.Name.IsExported() && (isGinContext(fn) || isEchoContext(fn)) {
				logrus.Infof("Found handler: %s from %s", fn.Name.Name, filename)
				handlers = append(handlers, fn)
			}
		}
	}
	return handlers, fset, nil
}

func isGinContext(fn *ast.FuncDecl) bool {
	if fn.Type.Params.NumFields() != 1 {
		return false
	}
	paramType := fn.Type.Params.List[0].Type

	if starExpr, ok := paramType.(*ast.StarExpr); ok {
		if selectorExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
			if selectorExpr.Sel.Name == "Context" {
				if ident, ok := selectorExpr.X.(*ast.Ident); ok {
					return ident.Name == "gin"
				}
			}
		}
	}
	return false
}

func isEchoContext(fn *ast.FuncDecl) bool {
	if fn.Type.Params.NumFields() != 1 {
		return false
	}
	paramType := fn.Type.Params.List[0].Type

	if selectorExpr, ok := paramType.(*ast.SelectorExpr); ok {
		if selectorExpr.Sel.Name == "Context" {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				return ident.Name == "echo"
			}
		}
	}
	return false
}
