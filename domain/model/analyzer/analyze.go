package analyzer

import (
	"fmt"
	"path/filepath"

	"github.com/lkeix/vis-gode-dependency/domain/model"
	"golang.org/x/tools/go/packages"
)

type Analizer struct {
}

func NewAnalizer() *Analizer {
	return &Analizer{}
}

var cfg = &packages.Config{
	Mode: packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedTypes |
		packages.NeedTypesInfo |
		packages.NeedName |
		packages.NeedImports |
		packages.NeedDeps,
}

func (a *Analizer) AnalyzeDependency() (model.DependencyList, error) {
	// files := make([]*model.File, 0)
	files, err := filepath.Glob("./go.mod")
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("go.mod file not found")
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, err
	}

	preAnalyzedPkgs, err := a.preAnalizePackages(pkgs)
	if err != nil {
		return nil, err
	}

	return preAnalyzedPkgs.Analyze()
}

func (a *Analizer) preAnalizePackages(pkgs []*packages.Package) (model.Packages, error) {
	ret := make(model.Packages, 0)
	for _, pkg := range pkgs {
		ret = append(ret, model.NewPackage(pkg.ID, pkg))
	}

	return ret, nil
}