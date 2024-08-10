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
	modName string
}

func NewPlantUML(modName string) *plantuml {
	return &plantuml{
		modName: modName,
	}
}

func (p *plantuml) Visualize(dependencyList *languagecomponents.DependencyList, outputPath string) error {
	plantUML := p.generateDiagram(dependencyList)

	fmt.Println(plantUML)

	return nil
}

func (p *plantuml) generateDiagram(dependencyList *languagecomponents.DependencyList) string {
	builder := strings.Builder{}

	builder.WriteString("@startuml\n")

	structDiagramBuilder := p.generatestructDiagram(dependencyList)
	builder.WriteString(structDiagramBuilder.String())

	dependencyDiagramBuilder := p.generateDependencyDiagram(dependencyList)
	builder.WriteString(dependencyDiagramBuilder.String())

	builder.WriteString("@enduml\n")

	return strings.ReplaceAll(strings.ReplaceAll(builder.String(), p.modName, ""), "/", ".")
}

func (p *plantuml) generatestructDiagram(dependencyList *languagecomponents.DependencyList) strings.Builder {
	pkgs := dependencyList.Aggregate()

	builder := strings.Builder{}

	for _, pkg := range pkgs {
		pkgDeclearation := fmt.Sprintf("package \"%s\" {\n", pkg.String())
		builder.WriteString(pkgDeclearation)
		for _, file := range pkg.Files {
			fileBaseName := filepath.Base(file.String())
			ext := filepath.Ext(fileBaseName)
			fileName := strings.Replace(fileBaseName, ext, "", 1)
			fileDeclearation := fmt.Sprintf("  package \"%s.%s\" {\n", pkg.String(), fileName)
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
				if object.ImplementInterface != nil {
					objectDeclearation = fmt.Sprintf("    struct %s implements \"%s\" {\n", object.Name, object.ImplementInterface.Name)

					if object.ImplementInterface.Package != nil {
						fileBaseName := filepath.Base(object.ImplementInterface.File.Name)
						ext := filepath.Ext(fileBaseName)
						fileName := strings.Replace(fileBaseName, ext, "", 1)
						objectDeclearation = fmt.Sprintf("    struct %s implements \"%s.%s.%s\" {\n", object.Name, object.ImplementInterface.Package.Name, fileName, object.ImplementInterface.Name)
					}
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
				if object.ImplementInterface != nil {
					objectDeclearation = fmt.Sprintf("    struct %s implements \"%s\" {\n", object.Name, object.ImplementInterface.Name)

					if object.ImplementInterface.Package != nil {
						fileBaseName := filepath.Base(object.ImplementInterface.File.Name)
						ext := filepath.Ext(fileBaseName)
						fileName := strings.Replace(fileBaseName, ext, "", 1)
						objectDeclearation = fmt.Sprintf("    struct %s implements \"%s.%s.%s\" {\n", object.Name, object.ImplementInterface.Package.Name, fileName, object.ImplementInterface.Name)
					}
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
		builder.WriteString("}\n\n")
	}

	return builder
}

func (p *plantuml) generateDependencyDiagram(dependencyList *languagecomponents.DependencyList) strings.Builder {
	builder := strings.Builder{}

	for _, dep := range dependencyList.List() {
		fromFileBaseName := filepath.Base(dep.FromFile.Name)
		ext := filepath.Ext(fromFileBaseName)
		fromFileName := strings.Replace(fromFileBaseName, ext, "", 1)

		toFileBaseName := filepath.Base(dep.ToFile.Name)
		toFileName := strings.Replace(toFileBaseName, ext, "", 1)
		writeDependency := fmt.Sprintf("%s.%s.%s ..> %s.%s.%s\n", dep.FromPackage.Name, fromFileName, dep.FromObject.Name, dep.ToPackage.Name, toFileName, dep.ToObject.Name)
		builder.WriteString(writeDependency)
	}

	return builder
}
