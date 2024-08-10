package languagecomponents_test

import (
	"testing"

	"github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"
)

func TestDependencyList_TopologicalSort(t *testing.T) {
	tests := []struct {
		name   string
		dl     languagecomponents.DependencyList
		answer languagecomponents.DependencyList
	}{
		{
			name:   "can sort dependency list",
			dl:     languagecomponents.DependencyList{},
			answer: languagecomponents.DependencyList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
