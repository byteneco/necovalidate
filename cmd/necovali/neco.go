package main

import (
	"github.com/byteneco/necovalidate"
)

func main() {
	p := necovalidate.NewParser()
	err := p.ParseFile("validate.go")
	if err != nil {
		panic(err)
	}
	p.Generate()
}
