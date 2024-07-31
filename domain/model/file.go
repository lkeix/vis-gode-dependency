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
			references := f.analyzeStatements(d.Body.List, pkgs)
			for _, reference := range references {
				obj := f.getObject(reference)
				fmt.Printf("%s.%s call %s.%s.%s\n", f.Name, d.Name.Name, reference.Pkg().Name(), obj.Name, reference.Id())
			}
		}
	}

	return nil
}

type reference struct {
	pkg *Package
	obj *Object
}

func (f *File) analyzeStatements(stmts []ast.Stmt, pkgs Packages) []types.Object {
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
					obj := pkg.pkg.TypesInfo.ObjectOf(selExpr.Sel)
					if obj != nil {
						ret = append(ret, obj)
					}
				}
			}
		}
	}

	return ret
}

func (f *File) getObject(obj types.Object) *Object {
	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return nil
	}

	if recv := sig.Recv(); recv != nil {
		recvType := recv.Type()
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}
		if named, ok := recvType.(*types.Named); ok {
			return NewObject(named.Obj().Name(), "struct", token.NoPos)
		}
	}

	return nil
}
