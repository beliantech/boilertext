package extractor

import (
	"bufio"
	"io"
	"strings"

	"github.com/PageDash/boilertext/pkg"
	"github.com/PageDash/boilertext/pkg/util"
)

// ShallowTextExtractor is an implementation of BoilerText
type ShallowTextExtractor struct {
	splitStrategy bufio.SplitFunc
}

// NewShallowTextExtractor returns a ShallowTextExtractor
func NewShallowTextExtractor(splitStrategy bufio.SplitFunc) *ShallowTextExtractor {
	return &ShallowTextExtractor{
		splitStrategy: splitStrategy,
	}
}

// Process takes raw HTML as an input and returns content text of that HTML minus the boilerplate.
func (s ShallowTextExtractor) Process(reader io.Reader) (string, error) {
	blocks, err := boilertext.GenerateTextBlocks(reader, s.splitStrategy)
	if err != nil {
		return "", err
	}

	// Block processing complete. Let's gather the wheat and discard the chaff.
	var contentText string
	var curr, prev, next *boilertext.TextBlock
	for i := range blocks {
		curr = blocks[i]
		util.Println("Block content", "NumOfWords", curr.NumOfWords, "NumOfAnchorWords", curr.NumOfAnchorWords, "Content", curr.Content)

		if i == 0 {
			prev = nil
		} else {
			prev = blocks[i-1]
		}
		if i == len(blocks)-1 {
			next = nil
		} else {
			next = blocks[i+1]
		}

		isContent := false
		if curr.LinkDensity() <= 0.333333 {
			if prev != nil && prev.LinkDensity() <= 0.555556 {
				if curr.NumOfWords <= 16 {
					if next != nil && next.NumOfWords <= 15 {
						if prev != nil && prev.NumOfWords <= 4 {
							isContent = false
						} else {
							isContent = true
						}
					} else {
						isContent = true
					}
				} else {
					isContent = true
				}
			} else {
				if curr.NumOfWords <= 40 {
					if next != nil && next.NumOfWords <= 17 {
						isContent = false
					} else {
						isContent = true
					}
				} else {
					isContent = true
				}
			}
		} else {
			isContent = false
		}

		if isContent {
			contentText += strings.TrimSpace(curr.Content) + " "
		}
	}

	return contentText, nil
}
