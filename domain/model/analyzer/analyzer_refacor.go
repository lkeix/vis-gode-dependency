package analyzer

import (
	"fmt"
	"path/filepath"

	"github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Analyzer struct {
	project *languagecomponents.Project
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze() (*languagecomponents.Project, error) {
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

	program, ssaPackages := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	program.Build()

	return languagecomponents.NewProject(program, ssaPackages, pkgs), err
}
