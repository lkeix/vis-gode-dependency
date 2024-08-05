package languagecomponents

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	Name        string
	Dir         string
	PackageName string
	Objects     []*Object
	Interfaces  []*Interface
	Funcs       []*Function

	syntax *ast.File
	fset   *token.FileSet
}

func NewFile(name string, syntax *ast.File, fset *token.FileSet) *File {
	funcs, interfaces := preAnalyze(syntax)

	return &File{
		Name:       name,
		syntax:     syntax,
		Funcs:      funcs,
		Interfaces: interfaces,
		fset:       fset,
	}
}

func preAnalyze(syntax *ast.File) ([]*Function, []*Interface) {
	retFuncs := make([]*Function, 0)
	retInteraces := make([]*Interface, 0)

	for _, decl := range syntax.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if i, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						methods := make([]*Function, 0, len(i.Methods.List))
						for _, field := range i.Methods.List {
							if _, ok := field.Type.(*ast.FuncType); ok {
								methods = append(methods, NewFunction(field.Names[0].Name, field.Pos()))
							}
						}
						retInteraces = append(retInteraces, NewInterface(typeSpec.Name.Name, methods, typeSpec.Pos()))
					}
				}
			}
		case *ast.FuncDecl:
			retFuncs = append(retFuncs, NewFunction(d.Name.Name, d.Pos()))
		}
	}

	return retFuncs, retInteraces
}

func (f *File) String() string {
	return f.Name
}

func (f *File) Analyze(pkgs Packages) (DependencyList, error) {
	deps := make(DependencyList, 0)
	for _, decl := range f.syntax.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			references := f.analyzeStatements(d.Body.List, pkgs)
			var fromObj *Object
			if d.Recv != nil && len(d.Recv.List) != 0 {
				recv := d.Recv.List[0]
				recvType := pkgs.FindReciverDeclarationByField(recv)
				if recvType != nil {
					objBaseName := filepath.Base(recvType.Type().String())
					replacePrefix := regexp.MustCompile(`[a-z]+.`).FindString(objBaseName)
					objName := strings.ReplaceAll(objBaseName, replacePrefix, "")
					fromObj = NewObject(objName, "struct", recv.Pos(), recvType)
				}
			}

			fromFun := NewFunction(d.Name.Name, d.Pos())

			for _, reference := range references {
				dep := NewDependency(f.lookupPackage(pkgs), f, fromObj, fromFun, reference.pkg, reference.obj.lookupFile(pkgs), reference.obj, reference.fun)
				deps = append(deps, dep)
			}
		}
	}

	return deps, nil
}

type reference struct {
	pkg *Package
	obj *Object
	fun *Function
}

func (f *File) analyzeStatements(stmts []ast.Stmt, pkgs Packages) []*reference {
	ret := make([]*reference, 0)
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.ExprStmt:
			switch e := s.X.(type) {
			case *ast.CallExpr:
				selExpr, ok := e.Fun.(*ast.SelectorExpr)
				if !ok {
					continue
				}

				fun := NewFunction(selExpr.Sel.Name, selExpr.Sel.Pos())
				for _, pkg := range pkgs {
					gotObj := pkg.pkg.TypesInfo.ObjectOf(selExpr.Sel)
					if gotObj != nil {
						obj := f.getObject(gotObj)
						p := pkgs.FindPackageByPath(obj.obj.Pkg().Path())

						ret = append(ret, &reference{
							pkg: p,
							obj: obj,
							fun: fun,
						})
					}
				}
			}
		}
	}

	return ret
}

func (f *File) getObject(obj types.Object) *Object {
	funcObj, ok := obj.(*types.Func)
	if !ok {
		return nil
	}

	sig, ok := funcObj.Type().(*types.Signature)
	if !ok {
		return nil
	}

	if recv := sig.Recv(); recv != nil {
		recvType := recv.Type()
		if ptr, ok := recvType.(*types.Pointer); ok {
			recvType = ptr.Elem()
		}
		if named, ok := recvType.(*types.Named); ok {
			return NewObject(named.Obj().Name(), "struct", recv.Pos(), obj)
		}
	}

	return nil
}

func (f *File) lookupPackage(pkgs Packages) *Package {
	for _, pkg := range pkgs {
		for i := range pkg.pkg.GoFiles {
			if pkg.pkg.Syntax[i] == f.syntax {
				return pkg
			}
		}
	}

	return nil
}

func (f *File) complete() {
	objMap := make(map[string]*Object)

	for _, decl := range f.syntax.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if i, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						interfaceMethods := make([]*Function, 0, len(i.Methods.List))
						for _, field := range i.Methods.List {
							if _, ok := field.Type.(*ast.FuncType); ok {
								interfaceMethods = append(interfaceMethods, NewFunction(field.Names[0].Name, field.Pos()))
							}
						}
						f.Interfaces = append(f.Interfaces, NewInterface(typeSpec.Name.Name, interfaceMethods, typeSpec.Pos()))
					}

					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						objMap[typeSpec.Name.Name] = NewObject(typeSpec.Name.Name, "struct", typeSpec.Pos(), nil)
					}
				}
			}
		case *ast.FuncDecl:
			if d.Recv != nil && len(d.Recv.List) != 0 {
				recvName := d.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
				obj := objMap[recvName]
				if obj != nil {
					obj.Methods = append(obj.Methods, NewFunction(d.Name.Name, d.Pos()))
					objMap[recvName] = obj
				}
			}
		}
	}

	for _, obj := range objMap {
		f.Objects = append(f.Objects, obj)
	}
}
