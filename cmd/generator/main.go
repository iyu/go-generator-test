package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/iyu/go-generator-test/cmd/generator/tpl"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("source required.")
	}
	source := os.Args[1]
	pkgName, entities, err := parse(source)
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New("repository.go").Parse(tpl.Repository)
	if err != nil {
		log.Fatal(err)
	}

	var w io.Writer
	if len(os.Args) >= 3 {
		outputPath := os.Args[2]
		file, err := os.Create(outputPath)
		if err != nil {
			log.Fatal(err)
		}
		w = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("Failed to close file %s: %+v\n", outputPath, err)
			}
		}()
	} else {
		w = os.Stdout
	}
	err = t.Execute(w, struct {
		PkgName  string
		Entities []*entity
	}{
		PkgName:  pkgName,
		Entities: entities,
	})
	if err != nil {
		log.Println(err)
	}
}

func parse(source string) (string, []*entity, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, source, nil, parser.Mode(0))
	if err != nil {
		return "", nil, err
	}
	entities := make([]*entity, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		if v, ok := n.(*ast.TypeSpec); ok {
			if t, ok := v.Type.(*ast.StructType); ok {
				fields := make([]*field, 0)
				for _, f := range t.Fields.List {
					if f.Tag == nil {
						continue
					}
					tagValue, err := strconv.Unquote(f.Tag.Value)
					if err != nil {
						panic(err)
					}
					tag := reflect.StructTag(tagValue)
					fields = append(fields, &field{
						Name:     f.Names[0].Name,
						JSONName: strings.Split(tag.Get("json"), ",")[0],
						Type:     f.Type.(*ast.Ident).Name,
					})
				}
				entities = append(entities, &entity{
					Name: v.Name.Name,
				})
			}
		}
		return true
	})
	return f.Name.Name, entities, nil
}

type entity struct {
	Name   string
	Fields []*field
}

type field struct {
	Name     string
	JSONName string
	Type     string
}
