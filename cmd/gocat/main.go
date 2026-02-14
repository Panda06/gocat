package main

import (
	"gocat/internal/argparser"
	"gocat/internal/cat"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	options, files, err := argparser.ParseArgs(args)
	if err != nil {
		log.Fatal(err)
	}
	cat.Run(options, files)
}
