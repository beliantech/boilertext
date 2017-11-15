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

	genBlocksBasedOnFlag := func(file *os.File) {
		defer file.Close()
		if *splitPtr == "word" {
			blocks, err = boilertext.GenerateTextBlocks(file, bufio.ScanWords)
			if err != nil {
				log.Fatal(err)
			}
		} else if *splitPtr == "rune" {
			blocks, err = boilertext.GenerateTextBlocks(file, bufio.ScanRunes)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Missing split argument")
		}
	}

	var ex boilertext.Extractor
	if *extractorPtr == "shallow" {
		ex = extractor.ShallowTextExtractor{}
		genBlocksBasedOnFlag(file)
	} else {
		// Returns all text
		ex = &extractor.AllTextExtractor{}
		genBlocksBasedOnFlag(file)

		// Calculate percentage of words that are links
		linkCount := 0
		linkWordCount := 0
		wordCount := 0
		for _, block := range blocks {
			linkWordCount += block.NumOfAnchorWords
			wordCount += block.NumOfWords

			if block.NumOfAnchorWords > 0 {
				linkCount++
			}
			fmt.Println("BLOCK", block)
		}

		fmt.Println("Percentage of words are links:", float64(linkWordCount)/float64(wordCount)*100.0, "%")
		fmt.Println("Percentage of blocks containing links:", float64(linkCount)/float64(len(blocks))*100.0, "%")
	}

	res, err := ex.Process(blocks)
	if err != nil {
		log.Fatal("Extractor failed")
	}

	fmt.Println("RESULT", res)
}
