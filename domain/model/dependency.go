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
