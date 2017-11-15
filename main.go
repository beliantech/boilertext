package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PageDash/boilertext/pkg/boilertext"
	"github.com/PageDash/boilertext/pkg/extractor"
)

func main() {
	splitPtr := flag.String("split", "word", "a string")
	extractorPtr := flag.String("extractor", "shallow", "a string")
	flag.Parse()

	file, err := os.Open("sample/" + os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal("Failed to open file")
	}

	var blocks []*boilertext.TextBlock

	var ex boilertext.Extractor
	if *extractorPtr == "shallow" {
		if *splitPtr == "word" {
			blocks, err = boilertext.GenerateTextBlocks(file, bufio.ScanWords)
			if err != nil {
				log.Fatal(err)
			}
			ex = extractor.ShallowTextExtractor{}
		} else if *splitPtr == "rune" {
			blocks, err = boilertext.GenerateTextBlocks(file, bufio.ScanRunes)
			if err != nil {
				log.Fatal(err)
			}
			ex = extractor.ShallowTextExtractor{}
		} else {
			log.Fatal("Missing split argument")
		}
	} else {
		// Returns all text
		blocks, err = boilertext.GenerateTextBlocks(file, bufio.ScanWords)

		// Calculate percentage of words that are links
		linkWordCount := 0
		wordCount := 0
		for _, block := range blocks {
			linkWordCount += block.NumOfAnchorWords
			wordCount += block.NumOfWords
		}
		fmt.Println("Percentage of words are links:", float64(linkWordCount)/float64(wordCount)*100.0, "%")

		ex = &extractor.AllTextExtractor{}
	}
	res, err := ex.Process(blocks)
	if err != nil {
		log.Fatal("Extractor failed")
	}

	fmt.Println("RESULT", res)
}
