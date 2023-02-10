package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	// "os"
	"time"

	"github.com/wonderstone/GEP-MOD/functions"
	"github.com/wonderstone/GEP-MOD/gene"
	"github.com/wonderstone/GEP-MOD/genome"
	"github.com/wonderstone/GEP-MOD/genomeset"
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

/*
Time Complexity: O(n^2)
*/
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

/*
Time Complexity: O(n*m)
*/
func validateFuncGS(gs *genomeset.GenomeSet) float64 {
	result := 0.0
	for _, n := range srTests {
		fitness := 0.0

		for _, g := range gs.Genomes {
			r := g.EvalMath(n.in)
			// fmt.Printf("r=%v, n.in=%v, n.out=%v, g=%v\n", r, n.in, n.out, g)
			if math.IsInf(r, 0) {
				return 0.0
			}
			fitness += r
			// fmt.Printf("r=%v, n.in=%v, n.out=%v, fitness=%v, g=%v\n", r, n.in, n.out, fitness, g)

		}

		fitness = math.Abs(fitness - n.out)

		fitness = 1000.0 / (1.0 + fitness) // fitness is normalized and max value is 1000

		result += fitness
	}
	return result / float64(len(srTests))
}

func main() {
	funcs := []gene.FuncWeight{
		{Symbol: "+", Weight: 1},
		{Symbol: "-", Weight: 1},
		{Symbol: "*", Weight: 1},
		// {Symbol: "/", Weight: 1},
	}

	// 打印华丽的分割线
	fmt.Println("=====================================================================================================")
	// 1. use genome mode to demo
	numIn := len(srTests[0].in)
	e := model.New(funcs, functions.Float64, 10, 4, 2, numIn, 0, "+", validateFunc)
	s := e.Evolve(2000, 50000, 0.8, 0.5, 3, 0.5, 3, 0.5, 0.01, 0.01, 0.01)
	tmp := validateFunc(s)
	fmt.Println("s score:", tmp)
	// Write out the Go source code for the solution.
	gr, err := grammars.LoadGoMathGrammar()
	if err != nil {
		log.Printf("unable to load grammar: %v", err)
	}
	helpers := make(grammars.HelperMap)

	// this part is the KES for genome, linked by "|"+g.LinkFunc+"|"
	fmt.Printf("// (a^4 + a^3 + a^2 + a) solution karva expression for Genome:\n// %q \n", s)
	// this part is the KE for one gene
	for _, g := range s.Genes {
		exp, _ := g.Expression(gr, helpers)
		fmt.Println(exp)
	}
	// 打印华丽的分割线
	fmt.Println("=====================================================================================================")
	// 2. use genomeset mode to demo
	es := model.NewGS(funcs, functions.Float64, 50, 4, 2, 2, numIn, 0, "+", validateFuncGS)
	ss := es.EvolveGS(1000, 50000, 0.8, 0.5, 3, 0.5, 3, 0.5, 0.01, 0.01, 0.01)
	fmt.Println("ss score:", ss.Score)

	for _, s := range ss.Genomes {

		fmt.Printf("// (a^4 + a^3 + a^2 + a) solution karva expression for GenomeSet inner Genome:\n// %q \n", s)
		for _, g := range s.Genes {
			exp, _ := g.Expression(gr, helpers)
			fmt.Println(exp)
		}

	}
}
