package languagecomponents

import (
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

type Project struct {
	program     *ssa.Program
	ssaPackages []*ssa.Package
	packages    []*packages.Package
}

func NewProject(program *ssa.Program, ssaPacakge []*ssa.Package, packages []*packages.Package) *Project {
	return &Project{
		program:     program,
		ssaPackages: ssaPacakge,
		packages:    packages,
	}
}

func (p *Project) Analyze() (Packages, *DependencyList, error) {
	return nil, nil, nil
}
