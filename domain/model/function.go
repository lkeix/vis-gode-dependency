package model

type Function struct {
	Name          string
	CallFunctions []*CallFunction
	CallStruct    []*Object
}

type CallFunction struct {
	Package  *Package
	Struct   *Object
	Function *Function
}
