package extractor

import (
	"bufio"
	"io"
	"strings"

	boilertext "github.com/PageDash/boilertext/pkg"
)

// AllTextExtractor returns all text in the document.
type AllTextExtractor struct{}

// Process takes raw HTML as an input and returns all text within that HTML.
func (a *AllTextExtractor) Process(reader io.Reader) (string, error) {
	blocks, err := boilertext.GenerateTextBlocks(reader, bufio.ScanWords)
	if err != nil {
		return "", err
	}

	var contentText string
	for _, block := range blocks {
		contentText += strings.TrimSpace(block.Content) + " "
	}

	return contentText, nil
}
