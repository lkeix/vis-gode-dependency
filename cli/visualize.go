package cli

import (
	"github.com/lkeix/vis-gode-dependency/analizer"
	"github.com/spf13/cobra"
)

func NewVisuazlize() *cobra.Command {

	visualizeCmd := &cobra.Command{
		Use:   "visualize",
		Short: "Visualize Go dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			a := analizer.NewAnalizer()
			a.AnalizeDependency()
		},
	}

	return visualizeCmd
}
