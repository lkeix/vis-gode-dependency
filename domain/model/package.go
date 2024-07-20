package model

type Package struct {
	Name  string
	Files []*File
}

type Packages []*Package
