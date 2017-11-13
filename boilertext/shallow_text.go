package boilertext

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TextBlock struct {
	NumOfWords           int
	NumOfAnchorTextWords int
	Content              []byte
}

func (t *TextBlock) LinkDensity() float64 {
	if t.NumOfWords != 0 && t.NumOfAnchorTextWords != 0 {
		return float64(t.NumOfAnchorTextWords) / float64(t.NumOfWords)
	}

	return 0.0
}

// ShallowTextExtractor is an implementation of BoilerText
type ShallowTextExtractor struct {
}

// Process takes raw HTML as an input and returns content text of that HTML minus the boilerplate.
func (s ShallowTextExtractor) Process(reader io.Reader) ([]byte, error) {
	node, err := html.Parse(reader)

	if err != nil {
		return nil, errors.Wrap(err, "Parse HTML error")
	}

	blocks := make([]*TextBlock, 0, 20)
	var bufferText bytes.Buffer
	var bufferAnchorText bytes.Buffer

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			fmt.Println("Text Node:", "Parent", n.Parent.DataAtom, "Data", n.Data, "NextSibling", n.NextSibling)
			switch n.Parent.DataAtom {
			case atom.A:
				bufferText.WriteString(strings.TrimSpace(n.Data))
				bufferAnchorText.WriteString(strings.TrimSpace(n.Data))
			case atom.Strike, atom.U, atom.B, atom.I, atom.Em, atom.Strong, atom.Span, atom.Sup, atom.Code, atom.Tt, atom.Sub, atom.Var, atom.Font:
				bufferText.WriteString(strings.TrimSpace(n.Data))
			case atom.Style, atom.Script, atom.Option, atom.Object, atom.Embed, atom.Applet, atom.Link, atom.Noscript:
				// Ignore
			default:
				// Generate a new block
				bufferText.WriteString(strings.TrimSpace(n.Data))

				// Retrieve bytes
				bufferTextBytes := bufferText.Bytes()

				textScanner := bufio.NewScanner(bytes.NewReader(bufferTextBytes))
				// Set the split function for the scanning operation.
				textScanner.Split(bufio.ScanWords)
				// Count the words.
				textCount := 0
				for textScanner.Scan() {
					textCount++
				}

				anchorTextScanner := bufio.NewScanner(&bufferAnchorText)
				// Set the split function for the scanning operation.
				anchorTextScanner.Split(bufio.ScanWords)
				// Count the words.
				anchorTextCount := 0
				for anchorTextScanner.Scan() {
					anchorTextCount++
				}

				blocks = append(blocks, &TextBlock{
					NumOfWords:           textCount,
					NumOfAnchorTextWords: anchorTextCount,
					Content:              bufferTextBytes,
				})

				// Reset buffers
				bufferText.Reset()
				bufferAnchorText.Reset()
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	// Block processing complete. Let's gather the wheat and discard the chaff.

	var contentText bytes.Buffer
	var prev, next *TextBlock
	for i, curr := range blocks {
		fmt.Println("Block content", string(curr.Content))

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
			contentText.Write(curr.Content)
			contentText.WriteString(" ")
		}
	}

	return contentText.Bytes(), nil
}
