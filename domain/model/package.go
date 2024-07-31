package model

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	Name  string
	Files []*File

	pkg *packages.Package
}

func NewPackage(name string, pkg *packages.Package) *Package {
	files := make([]*File, 0)
	for i, file := range pkg.GoFiles {
		f := NewFile(file, pkg.Syntax[i], pkg.Fset)
		files = append(files, f)
	}

	return &Package{
		Name:  name,
		Files: files,

		pkg: pkg,
	}
}

type Packages []*Package

func (p Packages) Analize() {
	for _, pkg := range p {
		for _, file := range pkg.Files {
			file.Analyze(p)
		}
	}
}

func (p Packages) FindReciverDeclaration(field *ast.Field) types.Type {
	for _, pkg := range p {
		t := pkg.pkg.TypesInfo.TypeOf(field.Type)
		if t != nil {
			return t
		}
	}
	return nil
}
