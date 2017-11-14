package boilertext

import "io"

// Extractor is an interface that processes incoming HTML and
// outputs text within HTML minus all the boilerplate
type Extractor interface {
	Process(html io.Reader) (string, error)
}
