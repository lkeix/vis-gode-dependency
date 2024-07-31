package model

type Dependency struct {
	FromPackage *Package
	FromObject  *Object
	FromFunc    *Function

	ToPackage *Package
	ToFunc    *Function
	ToObject  *Object
}

func NewDependency(fromPackage *Package, fromObject *Object, fromFunc *Function, toPackage *Package, toFunc *Function, toObject *Object) *Dependency {
	return &Dependency{
		FromPackage: fromPackage,
		FromObject:  fromObject,
		FromFunc:    fromFunc,
		ToPackage:   toPackage,
		ToFunc:      toFunc,
		ToObject:    toObject,
	}
}

type DependencyList []*Dependency
