package necovalidate

import (
	"errors"
	"fmt"
	"io"
)

var frameGenerators = map[string]Gen{
	Max: max,
}

type FieldInfo struct {
	ReceiverName string
	FiledType    string
	Name         string
	Args         []string
}

type Gen func(buffer io.Writer, info FieldInfo) error

func genValidateCodeFrame(buffer io.Writer, validateName string, info FieldInfo) error {
	gen, ok := frameGenerators[validateName]
	if !ok {
		return errors.New(fmt.Sprintf("there is no generator for validation: %s", validateName))
	}
	return gen(buffer, info)
}

func max(buffer io.Writer, info FieldInfo) error {
	_, err := fmt.Fprintf(buffer, "\tif len(%s.%s) > %s { return errors.New(\"false\") } ",
		info.ReceiverName,
		info.Name,
		info.Args[0],
	)
	return err
}
