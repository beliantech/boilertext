package boilertext

import (
	"bufio"
	"io"
	"strings"

	"github.com/PageDash/boilertext/logger"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Extractor is an interface that processes incoming HTML and
// outputs text within HTML minus all the boilerplate
type Extractor interface {
	Process(html io.Reader) (string, error)
}

// TextBlock represents a text block which may comprise of inline elements.
type TextBlock struct {
	NumOfWords       int
	NumOfAnchorWords int
	Content          string
}

// LinkDensity is the number of link text words divided by the total number of words in the block.
func (t *TextBlock) LinkDensity() float64 {
	if t.NumOfWords != 0 && t.NumOfAnchorWords != 0 {
		return float64(t.NumOfAnchorWords) / float64(t.NumOfWords)
	}

	return 0.0
}

// GenerateTextBlocks takes a reader containing HTML and generates a TextBlock array from it.
func GenerateTextBlocks(reader io.Reader, splitStrategy bufio.SplitFunc) ([]*TextBlock, error) {
	node, err := html.Parse(reader)
	if err != nil {
		return nil, errors.Wrap(err, "Parse HTML error")
	}

	blocks := make([]*TextBlock, 0, 20)
	var bufferText string
	var bufferAnchorText string

	var bufferAppend func(s string, isAnchor bool)
	bufferAppend = func(s string, isAnchor bool) {
		// Normalize whitespace to max 1 whitespace.
		// This does not preserve the original spacing in some cases, but we'd rather have extra spaces than joined words.
		s = strings.TrimSpace(s)
		if strings.HasSuffix(bufferText, " ") {
			bufferText += s
		} else {
			bufferText += " " + s
		}
		if isAnchor {
			if strings.HasSuffix(bufferAnchorText, " ") {
				bufferAnchorText += s
			} else {
				bufferAnchorText += " " + s
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
				if len(bufferText) == 0 {
					bufferText = ""
					bufferAnchorText = ""
					return
				}

				textScanner := bufio.NewScanner(strings.NewReader(bufferText))
				// Set the split function for the scanning operation.
				textScanner.Split(splitStrategy)
				// Count the words.
				textCount := 0
				for textScanner.Scan() {
					textCount++
				}

				anchorTextScanner := bufio.NewScanner(strings.NewReader(bufferAnchorText))
				// Set the split function for the scanning operation.
				anchorTextScanner.Split(splitStrategy)
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
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return blocks, nil
}
