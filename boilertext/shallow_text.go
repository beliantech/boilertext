package boilertext

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/PageDash/boilertext/logger"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TextBlock struct {
	NumOfWords       int
	NumOfAnchorWords int
	Content          []byte
}

func (t *TextBlock) LinkDensity() float64 {
	if t.NumOfWords != 0 && t.NumOfAnchorWords != 0 {
		return float64(t.NumOfAnchorWords) / float64(t.NumOfWords)
	}

	return 0.0
}

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
func (s ShallowTextExtractor) Process(reader io.Reader) ([]byte, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, errors.Wrap(err, "Parse HTML error")
	}

	blocks := make([]*TextBlock, 0, 20)
	var bufferText bytes.Buffer
	var bufferAnchorText bytes.Buffer

	var bufferAppend func(s string, isAnchor bool)
	bufferAppend = func(s string, isAnchor bool) {
		// Normalize whitespace to max 1 whitespace.
		// This does not preserve the original spacing in some cases, but we'd rather have extra spaces than joined words.
		s = strings.TrimSpace(s)
		if bytes.HasSuffix(bufferText.Bytes(), []byte(" ")) {
			bufferText.WriteString(s)
		} else {
			bufferText.WriteString(" ")
			bufferText.WriteString(s)
		}
		if isAnchor {
			if bytes.HasSuffix(bufferAnchorText.Bytes(), []byte(" ")) {
				bufferAnchorText.WriteString(s)
			} else {
				bufferAnchorText.WriteString(" ")
				bufferAnchorText.WriteString(s)
			}
		}
	}

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			trimmedData := strings.TrimSpace(n.Data)
			if trimmedData != "" {
				logger.Println("TEXT NODE", "Parent:", n.Parent.DataAtom, "Data:", n.Data, "NextSibling:", n.NextSibling)
			}

			switch n.Parent.DataAtom {
			case atom.A:
				if trimmedData != "" {
					logger.Println("ANCHOR", n.Data)
					bufferAppend(n.Data, true)
				}
			case atom.Strike, atom.U, atom.B, atom.I, atom.Em, atom.Strong, atom.Span, atom.Sup, atom.Code, atom.Tt, atom.Sub, atom.Var, atom.Font, atom.Time:
				// Don't append whitespace
				if trimmedData != "" {
					logger.Println("INLINE", n.Data)
					bufferAppend(n.Data, false)
				}
			case atom.Style, atom.Script, atom.Option, atom.Object, atom.Embed, atom.Applet, atom.Link, atom.Noscript:
				// Ignore
			default:
				// Generate a new block
				if trimmedData != "" {
					logger.Println("DEFAULT BLOCK DATA", n.Data)
					bufferAppend(n.Data, false)
				}

				// Retrieve bytes

				// Quit if nothing here
				if bufferText.Len() == 0 {
					bufferText.Reset()
					bufferAnchorText.Reset()
					return
				}

				bufferTextBytesCopy := make([]byte, bufferText.Len())
				copy(bufferTextBytesCopy, bufferText.Bytes())

				textScanner := bufio.NewScanner(bytes.NewReader(bufferTextBytesCopy))
				// Set the split function for the scanning operation.
				textScanner.Split(s.splitStrategy)
				// Count the words.
				textCount := 0
				for textScanner.Scan() {
					textCount++
				}

				anchorTextScanner := bufio.NewScanner(bytes.NewReader(bufferAnchorText.Bytes()))
				// Set the split function for the scanning operation.
				anchorTextScanner.Split(s.splitStrategy)
				// Count the words.
				anchorTextCount := 0
				for anchorTextScanner.Scan() {
					anchorTextCount++
				}

				blocks = append(blocks, &TextBlock{
					NumOfWords:       textCount,
					NumOfAnchorWords: anchorTextCount,
					Content:          bufferTextBytesCopy,
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
	var curr, prev, next *TextBlock
	for i := range blocks {
		curr = blocks[i]
		logger.Println("Block content", "NumOfWords", curr.NumOfWords, "NumOfAnchorWords", curr.NumOfAnchorWords, "Content", string(curr.Content))

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
			contentText.Write(bytes.TrimSpace(curr.Content))
			contentText.WriteString(" ")
		}
	}

	return contentText.Bytes(), nil
}
