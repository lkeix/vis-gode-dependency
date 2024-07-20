package model

type File struct {
	Name        string
	Dir         string
	PackageName string
	Structs     []*Object
	Funcs       []*Function
	Interfaces  []*Interface
}
