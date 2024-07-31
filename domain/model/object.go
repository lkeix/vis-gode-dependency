package model

import "go/token"

type Object struct {
	Name    string
	Type    string
	Pos     token.Pos
	Methods []*Function
}

func NewObject(name, t string, pos token.Pos) *Object {
	return &Object{
		Name: name,
		Type: t,
		Pos:  pos,
	}
}

func (o *Object) String() string {
	return o.Name
}
