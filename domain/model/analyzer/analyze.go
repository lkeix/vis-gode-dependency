package analyzer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"
	"golang.org/x/tools/go/packages"
)

type Analyzer struct {
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
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

func (a *Analyzer) AnalyzeDependency() (*languagecomponents.DependencyList, error) {
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

func (a *Analyzer) preAnalizePackages(pkgs []*packages.Package) (languagecomponents.Packages, error) {
	ret := make(languagecomponents.Packages, 0)
	for _, pkg := range pkgs {
		ret = append(ret, languagecomponents.NewPackage(pkg.ID, pkg))
	}

	for _, pkg := range ret {
		for _, file := range pkg.Files {
			ret = file.AnalyzeGenDecls(ret)
		}
	}

	return ret, nil
}

func (a *Analyzer) ModName() (string, error) {
	f, err := os.ReadFile("./go.mod")
	if err != nil {
		return "", err
	}

	ss := strings.Split(string(f), "\n")
	for _, s := range ss {
		if strings.HasPrefix(s, "module ") {
			return strings.ReplaceAll(s+"/", "module ", ""), nil
		}
	}

	return "", errors.New("failed to analyze go.mod")
}
