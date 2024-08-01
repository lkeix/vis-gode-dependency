package repository

import "github.com/lkeix/vis-gode-dependency/domain/model"

type Visualizer interface {
	Visualize(model.DependencyList) error
}
