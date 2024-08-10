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
			a := analyzer.NewAnalyzer()
			dependencyList, err := a.AnalyzeDependency()
			if err != nil {
				panic(err)
			}

			modName, err := a.ModName()
			if err != nil {
				panic(err)
			}

			sorted := dependencyList.TopologicalSort()

			visualizer := infrastructure.NewPlantUML(modName)
			if err := visualizer.Visualize(sorted); err != nil {
				panic(err)
			}
		},
	}

	return visualizeCmd
}
