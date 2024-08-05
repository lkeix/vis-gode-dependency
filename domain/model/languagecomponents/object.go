package languagecomponents

import (
	"go/token"
	"go/types"
)

type Interface struct {
	Name    string
	Methods []*Function
	Pos     token.Pos
}

func NewInterface(name string, methods []*Function, pos token.Pos) *Interface {
	return &Interface{
		Name:    name,
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
