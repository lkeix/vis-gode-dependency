package model

import "go/token"

type Function struct {
	Name string
	Pos  token.Pos
}

func NewFunction(name string, pos token.Pos) *Function {
	return &Function{
		Name: name,
		Pos:  pos,
	}
}

func (f *Function) String() string {
	return f.Name
}
