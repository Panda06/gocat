package main

import (
	"fmt"
	"gocat/internal/argparser"
	"gocat/internal/cat"
	"os"
)

func main() {
	args := os.Args[1:]
	options, files, err := argparser.ParseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	cat.Run(options, files)
}
