package languagecomponents

import (
	"fmt"
	"slices"
)

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

type DependencyList struct {
	list     []*Dependency
	packages Packages
}

func NewDependencyList(pkgs Packages) *DependencyList {
	return &DependencyList{
		list:     make([]*Dependency, 0),
		packages: pkgs,
	}
}

func (d *DependencyList) String() string {
	var s string
	for _, dep := range d.list {
		s += dep.String()
	}

	return s
}

func (d *DependencyList) Methods(obj *Object) []*Function {
	methods := make([]*Function, 0)
	methodMap := make(map[string]*Function)
	for _, dep := range d.list {
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
	for _, dep := range d.list {
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

func (d *DependencyList) Files(pkg *Package) []*File {
	files := make([]*File, 0)
	fileMap := make(map[string]*File)
	for _, dep := range d.list {
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

func (d *DependencyList) TopologicalSort() *DependencyList {
	graph := make(map[*Package]Packages)
	inDegree := make(map[*Package]int)

	for _, dep := range d.list {
		inDegree[dep.ToPackage] = 0
		inDegree[dep.FromPackage] = 0

		graph[dep.FromPackage] = append(graph[dep.FromPackage], dep.ToPackage)
		inDegree[dep.ToPackage]++
	}

	queue := make(Packages, 0)
	for pkg, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, pkg)
		}
	}

	ordered := make(Packages, 0)
	for len(queue) > 0 {
		pkg := queue[0]
		queue = queue[1:]

		ordered = append(ordered, pkg)

		for _, neighbor := range graph[pkg] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	slices.Reverse(ordered)
	ret := NewDependencyList(d.packages)
	for _, pkg := range ordered {
		for _, dep := range d.list {
			if dep.ToPackage == pkg {
				ret.list = append(ret.list, dep)
			}
		}
	}

	return ret
}

func (d *DependencyList) Aggregate() Packages {

	pkgs := make(Packages, 0)
	pkgMap := make(map[string]*Package)

	for _, dep := range d.list {
		if _, ok := pkgMap[dep.FromPackage.Name]; !ok {
			files := d.Files(dep.FromPackage)
			for _, file := range files {
				objects := d.Objects(file)
				for k, obj := range objects {
					methods := d.Methods(obj)
					objects[k].Methods = methods

					i := obj.lookupImplementInterface(dep.FromPackage)
					objects[k].ImplementInterface = i
				}
				file.Objects = append(file.Objects, objects...)
			}
			pkg := dep.FromPackage
			pkg.Files = files
			pkgs = append(pkgs, pkg)
			pkgMap[pkg.Name] = pkg
		}

		if _, ok := pkgMap[dep.ToPackage.Name]; !ok {
			dep.ToPackage.complete()
			files := d.Files(dep.ToPackage)
			pkg := dep.ToPackage
			pkg.Files = files
			pkgs = append(pkgs, pkg)
			for _, file := range files {
				objects := d.Objects(file)
				for k, obj := range objects {
					methods := d.Methods(obj)
					objects[k].Methods = methods
					p := obj.lookupImplementObjectPackage(d.packages)
					pkgs = append(pkgs, p)
					pkgMap[p.Name] = p
				}
				file.Objects = append(file.Objects, objects...)
			}
			pkgMap[pkg.Name] = pkg
		}
	}

	return pkgs.Unique()
}

func (d *DependencyList) List() []*Dependency {
	return d.list
}
