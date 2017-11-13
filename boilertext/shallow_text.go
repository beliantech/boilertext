package boilertext

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

// ShallowTextExtractor is an implementation of BoilerText
type ShallowTextExtractor struct {
}

// Process takes raw HTML as an input and returns content text of that HTML minus the boilerplate.
func (s ShallowTextExtractor) Process(reader io.Reader) ([]byte, error) {
	node, err := html.Parse(reader)

	if err != nil {
		return nil, errors.Wrap(err, "Parse HTML error")
	}

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			fmt.Println(n)
		} else if n.Type == html.TextNode {
			fmt.Println(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return []byte(""), nil
}
