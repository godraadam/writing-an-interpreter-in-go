package main

import (
	"bufio"
	"fmt"
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
	run(source)
}

func runRepl() {
	scanner := bufio.NewScanner(os.Stdin)
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
		run(source)

	}
}

func run(source string) {
	// do the work here
}
