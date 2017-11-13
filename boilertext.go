package boilertext

// BoilerText is an interface that processes incoming HTML and
// outputs text within HTML minus all the boilerplate
type BoilerText interface {
	Process(html []byte) ([]byte, error)
}
