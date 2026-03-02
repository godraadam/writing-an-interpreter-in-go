package main

import (
	"bufio"
	"fmt"
	"monkey/eval"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
)

func main() {
	var args = os.Args
	if len(args) > 2 {
		fmt.Println("Usage: monkey <script>")
		os.Exit(1)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runRepl()
	}
}

func runFile(pathToScript string) {
	bytes, err := os.ReadFile(pathToScript)
	if err != nil {
		fmt.Printf("Could not read file %s\n", pathToScript)
		os.Exit(1)
	}
	source := string(bytes[:])
	env := object.NewEnvironment()
	run(source, env, false)
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()
	for {
		fmt.Printf(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		source := scanner.Text()
		if source == "exit" {
			break
		}
		run(source, env, true)

	}
}

func run(source string, env *object.Environment, replMode bool) {
	// do the work here
	l := lexer.New(source)
	p := parser.New(l)
	ast := p.Parse()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	}
	// add env vars here I guess
	result := eval.Eval(ast, env)
	if replMode && result != nil {
		fmt.Println(result.Inspect())
	}
}
