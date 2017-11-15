package extractor

import (
	"bufio"
	"log"
	"os"
	"testing"
)

func BenchmarkShallowText(b *testing.B) {
	extractor := NewShallowTextExtractor(bufio.ScanWords)
	file, err := os.Open("../sample/nyt.html")
	if err != nil {
		log.Fatal("Failed to open file")
	}

	for n := 0; n < b.N; n++ {
		_, err := extractor.Process(file)
		if err != nil {
			log.Fatal("Extraction failed")
		}
	}
}
