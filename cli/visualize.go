package cli

import (
	"github.com/lkeix/vis-gode-dependency/analizer"
	"github.com/spf13/cobra"
)

func NewVisuazlize() *cobra.Command {
	var path string

	visualizeCmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize Go dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			a := analizer.NewAnalizer(path)
			a.AnalizeDependency()
		},
	}

	visualizeCmd.Flags().StringVarP(&path, "path", "", "Path to the project", ".")

	return visualizeCmd
}
