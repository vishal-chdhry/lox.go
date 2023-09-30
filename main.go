package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/vishal-chdhry/lox.go/scanner"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("usage: lox [stript]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		err := runPrompt()
		if err != nil {
			panic(err)
		}
	}
	os.Exit(0)
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = run(string(bytes))
	if err != nil {
		os.Exit(65)
	}
	return nil
}

func runPrompt() error {
	inputScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("> ")
		if !inputScanner.Scan() {
			break
		}

		ln := inputScanner.Text()
		run(ln)
	}
	return nil
}

func run(src string) error {
	scnr := scanner.NewScanner(src)
	tokens, err := scnr.ScanTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}
