package analizer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
		packages.NeedDeps,
}

func (a *Analizer) preAnalize() (model.Packages, error) {
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

	preAnalyzedPkgs.Analize()

	return model.Packages{}, nil
}

func (a *Analizer) AnalizeDependency() (model.Packages, error) {
	return a.preAnalize()
}

func (a *Analizer) preAnalizePackages(pkgs []*packages.Package) (model.Packages, error) {
	ret := make(model.Packages, 0)
	for _, pkg := range pkgs {
		ret = append(ret, model.NewPackage(pkg.ID, pkg))
	}

	return ret, nil
}

func (a *Analizer) analizeAST(path string) (*model.File, error) {
	fset := token.NewFileSet()
	fast, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, decl := range fast.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			// TODO: analize GenDecl
		case *ast.FuncDecl:
			if d.Recv != nil {
				// TODO: reference to reciver type
				a.referenceRecvType(d.Recv.List[0])
			}
			fmt.Println(d.Name.Name)
		}
	}

	return nil, nil
}

func (a *Analizer) referenceRecvType(field *ast.Field) *model.Object {
	switch d := field.Type.(type) {
	case *ast.Ident:
		return model.NewObject(d.Name, "defined type", d.Pos())
	case *ast.StarExpr:
		i, ok := d.X.(*ast.Ident)
		if !ok {
			return nil
		}
		return model.NewObject(i.Name, "pointer", d.Pos())
	}
	return nil
}
