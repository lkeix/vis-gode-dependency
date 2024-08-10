package languagecomponents

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/lkeix/vis-gode-dependency/utils"
)

type Interface struct {
	Name    string
	Package *Package
	File    *File
	Methods []*Function
	Pos     token.Pos
}

func NewInterface(name string, methods []*Function, pkg *Package, file *File, pos token.Pos) *Interface {
	return &Interface{
		Name:    name,
		Package: pkg,
		File:    file,
		Methods: methods,
		Pos:     pos,
	}
}

func (i *Interface) isImplemented(o *Object) bool {
	for _, method := range i.Methods {
		if !o.hasMethod(method) {
			return false
		}
	}

	return true
}

func (i *Interface) lookupImplementedObject(pkg *Package) *Object {
	for _, file := range pkg.Files {
		for _, obj := range file.Objects {
			if obj.lookupImplementInterface(pkg) == i {
				return obj
			}
		}
	}

	return nil
}

type Object struct {
	Name               string
	Type               string
	Pos                token.Pos
	Methods            []*Function
	ImplementInterface *Interface

	obj types.Object
}

func NewObject(name, t string, pos token.Pos, obj types.Object) *Object {
	return &Object{
		Name: name,
		Type: t,
		Pos:  pos,
		obj:  obj,
	}
}

func (o *Object) hasMethod(method *Function) bool {
	if !utils.IsUpperCase(method.Name) {
		return true
	}

	for _, m := range o.Methods {
		if m.Name == method.Name {
			return true
		}
	}

	return false
}

func (o *Object) String() string {
	return o.Name
}

func (o *Object) IsInterface() bool {
	return o.Type == "interface"
}

func (o *Object) IsStruct() bool {
	return o.Type == "struct"
}

func (o *Object) lookupFile(pkgs Packages) *File {
	sig := o.obj.Type().(*types.Signature)
	recv := sig.Recv()
	if recv == nil {
		return nil
	}

	for _, pkg := range pkgs {
		for i := range pkg.pkg.GoFiles {
			recvType := recv.Type()
			if ptr, ok := recvType.(*types.Pointer); ok {
				recvType = ptr.Elem()
			}

			if named, ok := recvType.(*types.Named); ok {
				structType := named.Obj()
				pos := pkg.pkg.Fset.Position(structType.Pos())
				return NewFile(pos.Filename, pkg.pkg.Syntax[i], pkg.pkg.Fset)
			}
		}
	}
	return nil
}

func (o *Object) lookupImplementInterface(pkg *Package) *Interface {
	interfaces := pkg.findInterfaces()

	for _, i := range interfaces {
		if i.isImplemented(o) {
			return i
		}
	}

	return nil
}

func (o *Object) lookupImplementObjectPackage(pkgs Packages) *Package {
	ret := make(Packages, 0)

	for _, pkg := range pkgs {
		if p, ok := pkg.pkg.Imports[o.obj.Pkg().Path()]; ok {
			for fi, file := range pkg.Files {
				if len(file.Objects) == 0 {
					file.complete(pkg)
				}

				for oi, obj := range file.Objects {
					isImplemented := true
					for _, method := range obj.Methods {
						if !o.hasMethod(method) {
							isImplemented = false
							break
						}
					}

					selector := getGenDeclPkgSelector(file)
					if selector != nil && obj.ImplementInterface == nil && isImplemented {
						i := NewInterface(o.Name, o.Methods, NewPackage(p.ID, p), file, o.Pos)
						obj.ImplementInterface = i
						pkg.Files[fi].Objects[oi] = obj
						return pkg
					}
				}
			}
		}
	}

	if len(ret) > 0 {
		for _, pkg := range ret {
			for _, file := range pkg.Files {
				for _, decl := range file.syntax.Decls {
					if gd, ok := decl.(*ast.GenDecl); ok {
						for _, spec := range gd.Specs {
							if vs, ok := spec.(*ast.ValueSpec); ok {
								if sel, ok := vs.Type.(*ast.SelectorExpr); ok {
									if x, ok := sel.X.(*ast.Ident); ok {
										if o.obj.Pkg().Name() == x.Name && sel.Sel.Name == o.Name {
											return pkg
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func getGenDeclPkgSelector(file *File) *ast.SelectorExpr {
	for _, decl := range file.syntax.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gd.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					if sel, ok := vs.Type.(*ast.SelectorExpr); ok {
						return sel
					}
				}
			}
		}
	}

	return nil
}
