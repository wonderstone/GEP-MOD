package main

import (
	"fmt"
	"testing"

	"github.com/wonderstone/GEP-MOD/functions"
	"github.com/wonderstone/GEP-MOD/gene"
)

func TestBuildExp(t *testing.T) {
	// give  a string to build the karva expression
	str := "+.d0.-.d0.c1.d2.d3.d3.d4"
	// give the functions.FuncType
	funcType := functions.Float64

	// get the gene from str
	g := gene.New(str, funcType)
	vars := g.GenerateMathFuncVars()
	fmt.Println(vars)

}
