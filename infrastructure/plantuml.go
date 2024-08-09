package infrastructure

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lkeix/vis-gode-dependency/domain/model/languagecomponents"
	"github.com/lkeix/vis-gode-dependency/domain/repository"
	"github.com/lkeix/vis-gode-dependency/utils"
)

var _ repository.Visualizer = &plantuml{}

type plantuml struct {
}

func NewPlantUML() *plantuml {
	return &plantuml{}
}

func (p *plantuml) Visualize(dependencyList *languagecomponents.DependencyList) error {
	plantUML := p.generatestructDiagram(dependencyList)

	fmt.Println(plantUML)

	return nil
}

func (p *plantuml) generatestructDiagram(dependencyList *languagecomponents.DependencyList) string {
	pkgs := dependencyList.Aggregate()

	builder := strings.Builder{}

	builder.WriteString("@startuml\n")

	for _, pkg := range pkgs {
		pkgDeclearation := fmt.Sprintf("package \"%s\" {\n", pkg.Name)
		builder.WriteString(pkgDeclearation)
		for _, file := range pkg.Files {
			fileBaseName := filepath.Base(file.String())
			ext := filepath.Ext(fileBaseName)
			fileName := strings.Replace(fileBaseName, ext, "", 1)
			fileDeclearation := fmt.Sprintf("  package \"%s.%s\" {\n", pkg.Name, fileName)
			builder.WriteString(fileDeclearation)
			for _, inf := range file.Interfaces {
				interfaceDeclearation := fmt.Sprintf("    interface %s {\n", inf.Name)
				builder.WriteString(interfaceDeclearation)
				for _, method := range inf.Methods {
					methodDeclearation := fmt.Sprintf("      %s()\n", method.Name)
					if utils.IsUpperCase(method.Name) {
						methodDeclearation = fmt.Sprintf("      +%s()\n", method.Name)
					}
					builder.WriteString(methodDeclearation)
				}
				builder.WriteString("    }\n")
			}

			for _, object := range file.Objects {
				if object.Type == "struct" {
					continue
				}

				objectDeclearation := fmt.Sprintf("    interface %s {\n", object.Name)
				if object.ImplementInterface != nil && object.ImplementInterface.Package != nil {
					fileBaseName := filepath.Base(object.ImplementInterface.File.Name)
					ext := filepath.Ext(fileBaseName)
					fileName := strings.Replace(fileBaseName, ext, "", 1)
					objectDeclearation = fmt.Sprintf("    struct %s implements \"%s.%s.%s\" {\n", object.Name, object.ImplementInterface.Package.Name, fileName, object.ImplementInterface.Name)
				}
				builder.WriteString(objectDeclearation)

				for _, method := range object.Methods {
					methodDeclearation := fmt.Sprintf("      %s()\n", method.Name)
					if utils.IsUpperCase(method.Name) {
						methodDeclearation = fmt.Sprintf("      +%s()\n", method.Name)
					}
					builder.WriteString(methodDeclearation)
				}

				builder.WriteString("    }\n")

			}

			for _, object := range file.Objects {
				if object.Type == "interface" {
					continue
				}

				objectDeclearation := fmt.Sprintf("    struct %s {\n", object.Name)
				if object.ImplementInterface != nil && object.ImplementInterface.Package != nil {
					fileBaseName := filepath.Base(object.ImplementInterface.File.Name)
					ext := filepath.Ext(fileBaseName)
					fileName := strings.Replace(fileBaseName, ext, "", 1)
					objectDeclearation = fmt.Sprintf("    struct %s implements \"%s.%s.%s\" {\n", object.Name, object.ImplementInterface.Package.Name, fileName, object.ImplementInterface.Name)
				}
				builder.WriteString(objectDeclearation)

				for _, method := range object.Methods {
					methodDeclearation := fmt.Sprintf("      %s()\n", method.Name)
					if utils.IsUpperCase(method.Name) {
						methodDeclearation = fmt.Sprintf("      +%s()\n", method.Name)
					}
					builder.WriteString(methodDeclearation)
				}

				builder.WriteString("    }\n")

			}
			builder.WriteString("  }\n")
		}
		builder.WriteString("}\n")
	}

	builder.WriteString("@enduml\n")

	ret := builder.String()

	return strings.ReplaceAll(ret, "github.com", "github_com")
}
