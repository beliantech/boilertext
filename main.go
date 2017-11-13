package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PageDash/boilertext/boilertext"
)

func main() {
	file, err := os.Open("sample/" + os.Args[1])
	if err != nil {
		log.Fatal("Failed to open file")
	}

	extractor := boilertext.ShallowTextExtractor{}
	res, err := extractor.Process(file)
	if err != nil {
		log.Fatal("Extractor failed")
	}

	fmt.Println("Result", string(res))
}
