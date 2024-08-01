package infrastructure

import (
	"fmt"

	"github.com/lkeix/vis-gode-dependency/domain/model"
	"github.com/lkeix/vis-gode-dependency/domain/repository"
)

var _ repository.Visualizer = &graphviz{}

type graphviz struct {
}

func NewGraphviz() *graphviz {
	return &graphviz{}
}

func (g *graphviz) Visualize(dependencyList model.DependencyList) error {
	pkgs := dependencyList.Aggregate()
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, object := range file.Objects {
				fmt.Println(object.Methods)
			}
		}
	}

	return nil
}
