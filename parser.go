package necovalidate

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(filename string) ([]StructSpec, error) {
	var structSpecs []StructSpec
	filetree, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return structSpecs, err
	}
	for _, del := range filetree.Decls {
		genDel, ok := del.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDel.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structSpec := p.parseStruct(ts.Name.Name, st)
			structSpecs = append(structSpecs, structSpec)
		}
	}
	return structSpecs, nil
}

type StructSpec struct {
	Name   string
	Fields []FieldSpec
}

type FieldSpec struct {
	Name        string
	Kind        string
	Expressions []string
}

func (p *Parser) parseStruct(name string, st *ast.StructType) StructSpec {
	ret := StructSpec{Name: name}
	if st == nil || st.Fields == nil {
		return ret
	}

	for _, field := range st.Fields.List {
		if field.Doc == nil {
			continue
		}
		var expressions []string
		for _, comment := range field.Doc.List {
			expressions = append(expressions, strings.TrimSpace(strings.TrimLeft(comment.Text, "/")))
		}
		ret.Fields = append(ret.Fields, FieldSpec{
			Name:        field.Names[0].Name,
			Expressions: expressions,
		})
	}

	return ret
}
