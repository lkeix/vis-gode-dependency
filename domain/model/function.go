package model

type Function struct {
	Name string
}

func NewFunction(name string) *Function {
	return &Function{
		Name: name,
	}
}

func (f *Function) String() string {
	return f.Name
}
