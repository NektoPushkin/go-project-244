package formatters

import (
	"errors"
)

type Node struct {
	K      string  `json:"key"`
	Status string  `json:"type"`
	OldVal any     `json:"value1,omitempty"`
	NewVal any     `json:"value2,omitempty"`
	OldSub bool    `json:"-"`
	NewSub bool    `json:"-"`
	Childs []*Node `json:"children,omitempty"`
}

const (
	StatusAdded     = "added"
	StatusRemoved   = "deleted"
	StatusChanged   = "changed"
	StatusUnchanged = "unchagned"
	StatusNested    = "nested"
)

type Formats map[string]func(diffTree []*Node) string

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
	formatters           = Formats{
		"stylish":   RenderStylish,
		"plain":     RenderPlain,
		"json":      RenderJSON,
		"jsonpatch": RenderJSONPatch,
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
