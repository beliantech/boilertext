package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PageDash/boilertext/pkg"
	"github.com/PageDash/boilertext/pkg/extractor"
)

func main() {
	splitPtr := flag.String("split", "word", "a string")
	flag.Parse()

	file, err := os.Open("sample/" + os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal("Failed to open file")
	}

	var ex boilertext.Extractor
	if *splitPtr == "word" {
		ex = extractor.NewShallowTextExtractor(bufio.ScanWords)
	} else if *splitPtr == "rune" {
		ex = extractor.NewShallowTextExtractor(bufio.ScanRunes)
	} else {
		log.Fatal("Missing split argument")
	}
	res, err := ex.Process(file)
	if err != nil {
		log.Fatal("Extractor failed")
	}

	fmt.Println("RESULT", res)
}
