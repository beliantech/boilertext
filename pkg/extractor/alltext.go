package extractor

import (
	"strings"

	"github.com/PageDash/boilertext/pkg/boilertext"
)

// AllTextExtractor returns all text in the document.
type AllTextExtractor struct{}

// Process takes raw HTML as an input and returns all text within that HTML.
func (a *AllTextExtractor) Process(blocks []*boilertext.TextBlock) (string, error) {
	var contentText string
	for _, block := range blocks {
		contentText += strings.TrimSpace(block.Content) + " "
	}

	return contentText, nil
}
