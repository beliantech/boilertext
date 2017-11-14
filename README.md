# BoilerText

BoilerText is a Go implementation of the algorithm to remove boilerplate text from HTML files as described by http://www.l3s.de/~kohlschuetter/boilerplate. The paper is found [here](http://www.l3s.de/~kohlschuetter/publications/wsdm187-kohlschuetter.pdf) (PDF). The intent of BoilerText output is for full-text search indexing.

The reference implementation is found in https://github.com/PageDash/boilerpipe (forked from https://github.com/kohlschutter/boilerpipe). This implementation does its best to mimick the algorithm described in the paper, but isn't 100% the same as the `boilerpipe` implementation.

By no means idiomatic Go. We'll get there. PRs welcome to clean up stuff or to add new algorithms.

## Language Support (Split Strategy)

There are two possible split strategies that you will want to consider. For English and English-like languages (which consists of words formed by a sequence of characters), the `bufio.ScanWords` `SplitFunc` is appropriate. For languages such as Chinese and Japanese (which consists of rune characters), use the `bufio.ScanRunes` `SplitFunc` to obtain the desired result. Obviously this is a simplistic view, but we gotta start somewhere.

Note that the research algorithm was based on the English language. YMMV for other languages. We found that replacing word split with rune split for runic languages performed decently.

See https://github.com/abadojack/whatlanggo for language detection feature support.
