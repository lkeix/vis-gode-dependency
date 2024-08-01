package model

import "fmt"

type Dependency struct {
	FromPackage *Package
	FromFile    *File
	FromObject  *Object
	FromFunc    *Function

	ToPackage *Package
	ToFile    *File
	ToObject  *Object
	ToFunc    *Function
}

const layout = `%s.%s.%s.%s -> %s.%s.%s.%s`

func NewDependency(fromPackage *Package, fromFile *File, fromObject *Object, fromFunc *Function, toPackage *Package, toFile *File, toObject *Object, toFunc *Function) *Dependency {
	return &Dependency{
		FromPackage: fromPackage,
		FromFile:    fromFile,
		FromObject:  fromObject,
		FromFunc:    fromFunc,
		ToPackage:   toPackage,
		ToFile:      toFile,
		ToObject:    toObject,
		ToFunc:      toFunc,
	}
}

func (d *Dependency) String() string {
	return fmt.Sprintf(`package: %s
file: %s
object: %s
function: %s
↓
package: %s
file: %s
object: %s
function: %s
`, d.FromPackage, d.FromFile, d.FromObject, d.FromFunc, d.ToPackage, d.ToFile, d.ToObject, d.ToFunc)
}

type DependencyList []*Dependency

func (d DependencyList) String() string {
	var s string
	for _, dep := range d {
		s += dep.String()
	}

	return s
}

func (d DependencyList) Methods(obj *Object) []*Function {
	methods := make([]*Function, 0)
	methodMap := make(map[string]*Function)
	for _, dep := range d {
		if dep.FromObject == obj {
			if _, ok := methodMap[dep.FromFunc.Name]; !ok {
				methodMap[dep.FromFunc.Name] = dep.FromFunc
				methods = append(methods, dep.FromFunc)
			}
		}

		if dep.ToObject == obj {
			if _, ok := methodMap[dep.ToFunc.Name]; !ok {
				methodMap[dep.ToFunc.Name] = dep.ToFunc
				methods = append(methods, dep.ToFunc)
			}
		}
	}

	return methods
}

func (d DependencyList) Objects(file *File) []*Object {
	objects := make([]*Object, 0)
	objectMap := make(map[string]*Object)
	for _, dep := range d {
		if dep.FromFile == file {
			if _, ok := objectMap[dep.FromObject.Name]; !ok {
				objectMap[dep.FromObject.Name] = dep.FromObject
				objects = append(objects, dep.FromObject)
			}
		}

		if dep.ToFile == file {
			if _, ok := objectMap[dep.ToObject.Name]; !ok {
				objectMap[dep.ToObject.Name] = dep.ToObject
				objects = append(objects, dep.ToObject)
			}
		}
	}

	return objects
}

func (d DependencyList) Files(pkg *Package) []*File {
	files := make([]*File, 0)
	fileMap := make(map[string]*File)
	for _, dep := range d {
		if dep.FromPackage == pkg {
			if _, ok := fileMap[dep.FromFile.Name]; !ok {
				fileMap[dep.FromFile.Name] = dep.FromFile
				files = append(files, dep.FromFile)
			}
		}

		if dep.ToPackage == pkg {
			if _, ok := fileMap[dep.ToFile.Name]; !ok {
				fileMap[dep.ToFile.Name] = dep.ToFile
				files = append(files, dep.ToFile)
			}
		}
	}

	return files
}

func (d DependencyList) Aggregate() Packages {
	pkgs := make(Packages, 0)
	pkgMap := make(map[string]*Package)

	for _, dep := range d {
		if _, ok := pkgMap[dep.FromPackage.Name]; !ok {
			files := d.Files(dep.FromPackage)
			for _, file := range files {
				objects := d.Objects(file)
				for k, obj := range objects {
					methods := d.Methods(obj)
					objects[k].Methods = methods
				}
				file.Objects = append(file.Objects, objects...)
			}
			pkg := dep.FromPackage
			pkg.Files = files
			pkgs = append(pkgs, pkg)
			pkgMap[pkg.Name] = pkg
		}

		if _, ok := pkgMap[dep.ToPackage.Name]; !ok {
			files := d.Files(dep.ToPackage)
			for _, file := range files {
				objects := d.Objects(file)
				for k, obj := range objects {
					methods := d.Methods(obj)
					objects[k].Methods = methods
				}
				file.Objects = append(file.Objects, objects...)
			}
			pkg := dep.ToPackage
			pkg.Files = files
			pkgs = append(pkgs, pkg)
			pkgMap[pkg.Name] = pkg
		}
	}

	return pkgs
}
