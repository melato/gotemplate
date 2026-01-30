package gotemplate

type PropertyParser func([]byte) (map[string]any, error)
type parsers map[string]PropertyParser

var propertyParsers parsers = make(parsers)

func SetParser(name string, parse PropertyParser) {
	propertyParsers[name] = parse
}

func init() {
	SetParser("json", ParseJson)
}
