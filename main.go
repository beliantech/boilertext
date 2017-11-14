package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PageDash/boilertext/boilertext"
)

func main() {
	splitPtr := flag.String("split", "word", "a string")
	flag.Parse()

	file, err := os.Open("sample/" + os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal("Failed to open file")
	}

	var extractor boilertext.Extractor
	if *splitPtr == "word" {
		extractor = boilertext.NewShallowTextExtractor(bufio.ScanWords)
	} else if *splitPtr == "rune" {
		extractor = boilertext.NewShallowTextExtractor(bufio.ScanRunes)
	} else {
		log.Fatal("Missing split argument")
	}
	res, err := extractor.Process(file)
	if err != nil {
		log.Fatal("Extractor failed")
	}

	fmt.Println("RESULT", res)
}
