package model

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

type File struct {
	Name        string
	Dir         string
	PackageName string
	Structs     []*Object
	Funcs       []*Function
	Interfaces  []*Interface

	syntax *ast.File
	fset   *token.FileSet
}

func NewFile(name string, syntax *ast.File, fset *token.FileSet) *File {
	return &File{
		Name:   name,
		syntax: syntax,
		fset:   fset,
	}
}

func (f *File) Analyze(pkgs Packages) error {
	for _, decl := range f.syntax.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:

		case *ast.FuncDecl:
			references := f.analizeStatements(d.Body.List, pkgs)
			for _, reference := range references {
				fmt.Printf("%s.%s call %s\n", f.Name, d.Name.Name, reference)
			}
		}
	}

	return nil
}

func (f *File) analizeStatements(stmts []ast.Stmt, pkgs Packages) []types.Object {
	ret := make([]types.Object, 0)
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			switch e := s.X.(type) {
			case *ast.CallExpr:
				selExpr, ok := e.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				for _, pkg := range pkgs {
					if obj := pkg.pkg.Types.Scope().Lookup(selExpr.Sel.Name); obj != nil {
						ret = append(ret, obj)
					}
				}
			}
		}
	}
	return ret
}
