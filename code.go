package code

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"

	"github.com/fatih/color"

	"code/parsers"
)

var (
	ErrUnsupportedFormat = errors.New("unknown format, expected \"stylish\" or \"plain\"")
	noDifferenceFound    = "no differences found"
)

type diff map[string][2]any

func GenDiff(filepath1, filepath2, format string) (string, error) {
	config1, err := parsers.ParseFile(filepath1)

	if err != nil {
		return "", err
	}

	config2, err := parsers.ParseFile(filepath2)

	if err != nil {
		return "", err
	}

	if reflect.DeepEqual(config1, config2) {
		return noDifferenceFound, nil
	}

	diff := makeDiff(config1, config2)

	var ret string
	switch format {
	case "stylish":
		ret = renderStylish(diff)
	case "plain":
		ret = renderPlain(diff)
	default:
		return "", ErrUnsupportedFormat
	}

	return ret, nil
}

func makeDiff(config1, config2 map[string]any) diff {
	keys := make([]string, 0, len(config1)+len(config2))

	for k := range config1 {
		keys = append(keys, k)
	}

	for k := range config2 {
		if _, ok := config1[k]; !ok {
			keys = append(keys, k)
		}
	}

	diff := make(diff, len(keys))
	for _, k := range keys {
		diff[k] = [2]any{config1[k], config2[k]}
	}

	return diff
}

func renderStylish(diff diff) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for _, k := range slices.Sorted(maps.Keys(diff)) {
		oldVal, newVal := diff[k][0], diff[k][1]
		switch {
		case oldVal != nil && newVal == nil:
			fmt.Fprintf(&sb, "  %s %s: %v\n", color.RedString("-"), k, oldVal)
		case oldVal == nil && newVal != nil:
			fmt.Fprintf(&sb, "  %s %s: %v\n", color.GreenString("+"), k, newVal)
		case oldVal != newVal:
			if oldVal != newVal {
				fmt.Fprintf(&sb, "  %s %s: %v\n", color.RedString("-"), k, oldVal)
				fmt.Fprintf(&sb, "  %s %s: %v\n", color.GreenString("+"), k, newVal)
			}
		}
	}

	sb.WriteString("}")

	return sb.String()
}

func renderPlain(diff diff) string {
	return ""
}
