package model

import (
	"go/token"
	"go/types"
)

type Object struct {
	Name    string
	Type    string
	Pos     token.Pos
	Methods []*Function

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
