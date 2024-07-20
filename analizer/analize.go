package analizer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/lkeix/vis-gode-dependency/domain/model"
)

type Analizer struct {
	path string
}

func NewAnalizer(path string) *Analizer {
	return &Analizer{
		path: path,
	}
}

func (a *Analizer) preAnalize() model.Packages {
	files := make([]*model.File, 0, 0)

	filepath.Walk(a.path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".go" {
			return nil
		}

		f, err := a.analizeAST(path)
		if err != nil {
			return err
		}

		files = append(files, f)
		return nil
	})
	return model.Packages{}
}

func (a *Analizer) AnalizeDependency() model.Packages {
	return a.preAnalize()
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
