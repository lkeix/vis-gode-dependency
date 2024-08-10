package cli

import (
	"fmt"

	"github.com/lkeix/vis-gode-dependency/domain/model/analyzer"
	"github.com/lkeix/vis-gode-dependency/infrastructure"
	"github.com/spf13/cobra"
)

func NewVisuazlize() *cobra.Command {
	var output string
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

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("vizualize error: %v\n", err)
				}

				if err := visualizer.Cleanup(); err != nil {
					panic(err)
				}
			}()

			if err := visualizer.Visualize(sorted, output); err != nil {
				panic(err)
			}
		},
	}

	visualizeCmd.Flags().StringVarP(&output, "output", "o", "output", "specify plantuml diagram")

	return visualizeCmd
}
