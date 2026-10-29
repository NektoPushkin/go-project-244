package formatters

import (
	"errors"
)

type Node struct {
	K              string
	Status         int
	OldVal, NewVal any
	OldSub, NewSub bool
	Childs         []*Node
}

const (
	StatusAdded int = iota
	StatusRemoved
	StatusChanged
	StatusUnchanged
)

type Formats map[string]func(diffTree []*Node) string

var (
	ErrUnsupportedFormat = errors.New("unknown format, expected \"stylish\" or \"plain\"")
	formatters           = Formats{
		"stylish": RenderStylish,
		"plain":   RenderPlain,
	}
)

func Render(diffTree []*Node, format string) (string, error) {
	if formatterFunc, ok := formatters[format]; ok {
		return formatterFunc(diffTree), nil
	}
	return "", ErrUnsupportedFormat
}

func IsFormatAvailable(format string) bool {
	_, ok := formatters[format]
	return ok
}
