package main

import (
	"fmt"
	"github.com/byteneco/necovalidate"
)

func main() {
	spec, err := necovalidate.NewParser().ParseFile("validate.go")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", spec)
}
