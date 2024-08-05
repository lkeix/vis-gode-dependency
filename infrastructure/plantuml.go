package infrastructure

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"
	"github.com/lkeix/vis-gode-dependency/domain/repository"
)

var _ repository.Visualizer = &plantuml{}

type plantuml struct {
}

func NewPlantUML() *plantuml {
	return &plantuml{}
}

func (p *plantuml) Visualize(dependencyList languagecomponents.DependencyList) error {
	plantUML := p.generateClassDiagram(dependencyList)

	fmt.Println(plantUML)

	return nil
}

func (p *plantuml) generateClassDiagram(dependencyList languagecomponents.DependencyList) string {
	pkgs := dependencyList.Aggregate()

	builder := strings.Builder{}

	builder.WriteString("@startuml\n")

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			fileBaseName := filepath.Base(file.String())
			ext := filepath.Ext(fileBaseName)
			fileName := strings.Replace(fileBaseName, ext, "", 1)
			fileDeclearation := fmt.Sprintf("package \"%s.%s\" {\n", pkg.Name, fileName)
			builder.WriteString(fileDeclearation)
			for _, object := range file.Objects {
				i := object.ImplementInterface
				if i != nil {
					interfaceDeclearation := fmt.Sprintf("  interface %s {\n", i.Name)
					builder.WriteString(interfaceDeclearation)
					for _, method := range i.Methods {
						methodDeclearation := fmt.Sprintf("    %s()\n", method.Name)
						if isUpperCase(method.Name) {
							methodDeclearation = fmt.Sprintf("   +%s()\n", method.Name)
						}
						builder.WriteString(methodDeclearation)
					}
					builder.WriteString("  }\n")
				}

				objectDeclearation := fmt.Sprintf("  class %s {\n", object.Name)
				if object.ImplementInterface != nil {
					objectDeclearation = fmt.Sprintf("  class %s implements %s {\n", object.Name, object.ImplementInterface.Name)
				}
				builder.WriteString(objectDeclearation)

				for _, method := range object.Methods {
					methodDeclearation := fmt.Sprintf("    %s()\n", method.Name)
					if isUpperCase(method.Name) {
						methodDeclearation = fmt.Sprintf("   +%s()\n", method.Name)
					}
					builder.WriteString(methodDeclearation)
				}

				builder.WriteString("  }\n")

			}

			builder.WriteString("}\n")
		}
	}
	builder.WriteString("@enduml\n")

	ret := builder.String()

	return strings.ReplaceAll(ret, "github.com", "github_com")
}

func isUpperCase(s string) bool {
	return strings.ToUpper(string(s[0])) == string(s[0])
}
