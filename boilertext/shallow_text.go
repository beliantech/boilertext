package boilertext

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TextBlock struct {
	NumOfWords       int
	NumOfAnchorWords int
	Content          string
}

func (t *TextBlock) LinkDensity() float64 {
	if t.NumOfWords != 0 && t.NumOfAnchorWords != 0 {
		return float64(t.NumOfAnchorWords) / float64(t.NumOfWords)
	}

	return 0.0
}

// ShallowTextExtractor is an implementation of BoilerText
type ShallowTextExtractor struct {
}

// Process takes raw HTML as an input and returns content text of that HTML minus the boilerplate.
func (s ShallowTextExtractor) Process(reader io.Reader) (string, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return "", errors.Wrap(err, "Parse HTML error")
	}

	blocks := make([]*TextBlock, 0, 20)
	var bufferText string
	var bufferAnchorText string

	prevWasInline := false
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			trimmedData := strings.TrimSpace(n.Data)
			if trimmedData != "" {
				fmt.Println("TEXT NODE", "Parent:", n.Parent.DataAtom, "Data:", n.Data, "NextSibling:", n.NextSibling)
			}

			switch n.Parent.DataAtom {
			case atom.A:
				if trimmedData != "" {
					fmt.Println("ANCHOR", n.Data)
					if strings.HasSuffix(bufferText, " ") {
						bufferText += n.Data
					} else {
						bufferText += " " + n.Data
					}

					if strings.HasSuffix(bufferAnchorText, " ") {
						bufferAnchorText += n.Data
					} else {
						bufferAnchorText += " " + n.Data
					}
				}
				prevWasInline = true
			case atom.Strike, atom.U, atom.B, atom.I, atom.Em, atom.Strong, atom.Span, atom.Sup, atom.Code, atom.Tt, atom.Sub, atom.Var, atom.Font, atom.Time:
				// Don't append whitespace
				if trimmedData != "" {
					fmt.Println("INLINE", n.Data)
					if strings.HasSuffix(bufferText, " ") {
						bufferText += n.Data
					} else {
						bufferText += " " + n.Data
					}
				}
				prevWasInline = true
			case atom.Style, atom.Script, atom.Option, atom.Object, atom.Embed, atom.Applet, atom.Link, atom.Noscript:
				// Ignore
			default:
				// Generate a new block
				if trimmedData != "" {
					if strings.HasSuffix(bufferText, " ") {
						bufferText += strings.TrimLeftFunc(n.Data, unicode.IsSpace)
					} else {
						if prevWasInline {
							// No whitespace prepend when previous text node was inline block.
							bufferText += n.Data
						} else {
							bufferText += " " + strings.TrimLeftFunc(n.Data, unicode.IsSpace)
						}
					}
					fmt.Println("DEFAULT BLOCK DATA", n.Data)
				}

				// Retrieve bytes

				// Quit if nothing here
				if len(bufferText) == 0 {
					bufferText = ""
					bufferAnchorText = ""
					return
				}

				textScanner := bufio.NewScanner(strings.NewReader(bufferText))
				// Set the split function for the scanning operation.
				textScanner.Split(bufio.ScanWords)
				// Count the words.
				textCount := 0
				for textScanner.Scan() {
					textCount++
				}

				anchorTextScanner := bufio.NewScanner(strings.NewReader(bufferAnchorText))
				// Set the split function for the scanning operation.
				anchorTextScanner.Split(bufio.ScanWords)
				// Count the words.
				anchorTextCount := 0
				for anchorTextScanner.Scan() {
					anchorTextCount++
				}

				blocks = append(blocks, &TextBlock{
					NumOfWords:       textCount,
					NumOfAnchorWords: anchorTextCount,
					Content:          bufferText,
				})

				// Reset buffers
				bufferText = ""
				bufferAnchorText = ""
				prevWasInline = false
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	// Block processing complete. Let's gather the wheat and discard the chaff.

	var contentText string
	var curr, prev, next *TextBlock
	for i := range blocks {
		curr = blocks[i]
		fmt.Println("Block content", "NumOfWords", curr.NumOfWords, "NumOfAnchorWords", curr.NumOfAnchorWords, "Content", string(curr.Content))

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
