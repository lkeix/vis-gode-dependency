package main

import "github.com/lkeix/vis-gode-dependency/cli"

func main() {
	cmd := cli.NewCLI()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
