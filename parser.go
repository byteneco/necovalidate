package necovalidate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Parser struct {
	parsedFiles []ParsedFile
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(filename string) error {
	var structSpecs []StructSpec
	astFile, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, del := range astFile.Decls {
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
			structSpec, err := p.parseStruct(ts.Name.Name, st)
			if err != nil {
				return err
			}
			structSpecs = append(structSpecs, structSpec)
		}
	}
	p.parsedFiles = append(p.parsedFiles, ParsedFile{
		FileName:    filename,
		FileAst:     astFile,
		StructSpecs: structSpecs,
	})
	return nil

}

func (p *Parser) parseStruct(name string, st *ast.StructType) (StructSpec, error) {
	ret := StructSpec{Name: name}
	if st == nil || st.Fields == nil {
		return ret, nil
	}

	for _, field := range st.Fields.List {
		if field.Doc == nil {
			continue
		}
		var expressions []Expression
		for _, comment := range field.Doc.List {
			if isExpComment(comment.Text) {
				expression, err := parseComment(comment.Text)
				if err != nil {
					return ret, err
				}
				expressions = append(expressions, expression)
			}
		}
		ret.Fields = append(ret.Fields, FieldSpec{
			Name:        field.Names[0].Name,
			Expressions: expressions,
		})
	}

	return ret, nil
}

func (p *Parser) Generate() {
	for _, f := range p.parsedFiles {
		newFileGen(f).Generate()
	}
}

type fileGen struct {
	file              ParsedFile
	buffer            strings.Builder
	curStruct         StructSpec
	curStructReceiver string
}

func newFileGen(file ParsedFile) *fileGen {
	return &fileGen{
		file: file,
	}
}

func (f *fileGen) Generate() error {
	f.writePackage()
	for _, st := range f.file.StructSpecs {
		f.curStructReceiver = string([]rune(strings.ToLower(st.Name))[0])
		f.curStruct = st
		var err error
		fmt.Fprintf(&f.buffer, "\nfunc (%s %s) validate() error {\n", f.curStructReceiver, st.Name)
		f.GenStructValidation()
		fmt.Fprintf(&f.buffer, "\n}\n")
		if err != nil {
			return err
		}
	}
	wf, err := os.Create(f.genFileName())
	if err != nil {
		return err
	}
	wf.WriteString(f.buffer.String())
	return nil
}

func (f *fileGen) writePackage() {
	fmt.Fprintf(&f.buffer, "package %s\n", f.file.FileAst.Name.String())
}

func (f *fileGen) GenStructValidation() {
	for _, field := range f.curStruct.Fields {
		for _, exp := range field.Expressions {
			genValidateCodeFrame(&f.buffer, exp.Name, FieldInfo{
				ReceiverName: f.curStructReceiver,
				FiledType:    field.Kind,
				Name:         field.Name,
				Args:         exp.Args,
			})
		}
	}
}

func (f *fileGen) genFileName() string {
	name := strings.TrimSuffix(f.file.FileName, ".go")
	return fmt.Sprintf("%s_validation.go", name)
}
