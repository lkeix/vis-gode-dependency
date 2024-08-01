package cli

import (
	"github.com/lkeix/vis-gode-dependency/domain/model/analyzer"
	"github.com/lkeix/vis-gode-dependency/infrastructure"
	"github.com/spf13/cobra"
)

func NewVisuazlize() *cobra.Command {

	visualizeCmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize Go dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			a := analyzer.NewAnalizer()
			dependencyList, err := a.AnalyzeDependency()
			if err != nil {
				panic(err)
			}

			visualizer := infrastructure.NewGraphviz()
			if err := visualizer.Visualize(dependencyList); err != nil {
				panic(err)
			}
		},
	}

	return visualizeCmd
}
