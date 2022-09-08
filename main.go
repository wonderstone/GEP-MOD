package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/wonderstone/GEP-MOD/functions"
	"github.com/wonderstone/GEP-MOD/gene"
	"github.com/wonderstone/GEP-MOD/genome"
	"github.com/wonderstone/GEP-MOD/grammars"
	"github.com/wonderstone/GEP-MOD/model"
)

// srTests is a random sample of inputs and outputs for the function "a^4 + a^3 + a^2 + a"
var srTests = []struct {
	in  []float64
	out float64
}{
	// {[]float64{0}, 0},
	{[]float64{2.81}, 95.2425},
	{[]float64{6}, 1554},
	{[]float64{7.043}, 2866.55},
	{[]float64{8}, 4680},
	{[]float64{10}, 11110},
	{[]float64{11.38}, 18386},
	{[]float64{12}, 22620},
	{[]float64{14}, 41370},
	{[]float64{15}, 54240},
	{[]float64{20}, 168420},
	{[]float64{100}, 101010100},
	{[]float64{-100}, 99009900},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func validateFunc(g *genome.Genome) float64 {
	result := 0.0
	for _, n := range srTests {
		r := g.EvalMath(n.in)
		// fmt.Printf("r=%v, n.in=%v, n.out=%v, g=%v\n", r, n.in, n.out, g)
		if math.IsInf(r, 0) {
			return 0.0
		}
		fitness := math.Abs(r - n.out)
		fitness = 1000.0 / (1.0 + fitness) // fitness is normalized and max value is 1000
		// fmt.Printf("r=%v, n.in=%v, n.out=%v, fitness=%v, g=%v\n", r, n.in, n.out, fitness, g)
		result += fitness
	}
	return result / float64(len(srTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
		{Symbol: "/", Weight: 1},
	}
	numIn := len(srTests[0].in)
	e := model.New(funcs, functions.Float64, 50, 10, 2, numIn, 0, "+", validateFunc)
	s := e.Evolve(10000, 50000, 0.8, 0.5, 3, 0.5, 3, 0.5, 0.01, 0.01, 0.01)

	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoMathGrammar()
	if err != nil {
		log.Printf("unable to load grammar: %v", err)
	}
	fmt.Printf("\n// gepModel is auto-generated Go source code for the\n")
	fmt.Printf("// (a^4 + a^3 + a^2 + a) solution karva expression:\n// %q, score=%v\n", s, validateFunc(s))
	s.Write(os.Stdout, gr)
}
