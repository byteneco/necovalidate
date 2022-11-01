package necovalidate

import "go/ast"

type ParsedFile struct {
	FileName    string
	FileAst     *ast.File
	StructSpecs []StructSpec
}

type StructSpec struct {
	Name   string
	Fields []FieldSpec
}

type FieldSpec struct {
	Name        string
	Kind        string
	Expressions []Expression
}

type Expression struct {
	Name string
	Args []string
}
