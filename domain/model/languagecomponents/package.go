package languagecomponents

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
		pkg:   pkg,
	}
}

func (p *Package) findInterfaces() []*Interface {
	ret := make([]*Interface, 0)
	for _, file := range p.Files {
		ret = append(ret, file.Interfaces...)
	}

	return ret
}

func (p *Package) complete() {
	for i := range p.pkg.GoFiles {
		p.Files[i].complete(p)
	}
}

func (p *Package) String() string {
	return p.Name
}

type Packages []*Package

func (p Packages) Analyze() (*DependencyList, error) {
	dependencyList := NewDependencyList(p)
	for _, pkg := range p {
		for _, file := range pkg.Files {
			d, err := file.Analyze(p)
			if err != nil {
				return nil, err
			}

			dependencyList.list = append(dependencyList.list, d.list...)
		}

		pkg.findInterfaces()
	}

	return dependencyList, nil
}

func (p Packages) FindReciverDeclarationByField(field *ast.Field) types.Object {
	for _, pkg := range p {
		t := pkg.pkg.TypesInfo.ObjectOf(field.Names[0])
		if t != nil {
			return t
		}
	}
	return nil
}

func (p Packages) FindPackageByPath(path string) *Package {
	for _, pkg := range p {
		if pkg.pkg.PkgPath == path {
			return pkg
		}
	}

	return nil
}

func (p Packages) Packages() []*packages.Package {
	ret := make([]*packages.Package, 0)
	for _, pkg := range p {
		ret = append(ret, pkg.pkg)
	}

	return ret
}

func (p Packages) Unique() Packages {
	mp := make(map[*Package]struct{})
	ret := make(Packages, 0, len(p))

	for _, pkg := range p {
		if _, ok := mp[pkg]; !ok {
			ret = append(ret, pkg)
		}
		mp[pkg] = struct{}{}
	}

	return ret
}
