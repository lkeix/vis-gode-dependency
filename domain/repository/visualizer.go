package repository

import "github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"

type Visualizer interface {
	Visualize(dependenecyList *languagecomponents.DependencyList, output string) error
}
