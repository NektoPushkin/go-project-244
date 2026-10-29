package code

import (
	"errors"
	"fmt"
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

const (
	sAdded int = iota
	sRemoved
	sChanged
	sUnchanged
	sParent
)

type node struct {
	k              string
	status         int
	oldVal, newVal any
	oldSub, newSub bool
	childs         []*node
}

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

	diff := buildDiffTree(config1, config2)

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

func buildDiffTree(config1, config2 map[string]any) []*node {
	keys := make([]string, 0, len(config1)+len(config2))

	for k := range config1 {
		keys = append(keys, k)
	}

	for k := range config2 {
		if _, ok := config1[k]; !ok {
			keys = append(keys, k)
		}
	}
	slices.Sort(keys)
	nodes := make([]*node, 0, len(keys))

	for _, k := range keys {
		v1, ok1 := config1[k]
		v2, ok2 := config2[k]

		m1, okv1 := v1.(map[string]any)
		m2, okv2 := v2.(map[string]any)

		node := &node{k: k, childs: buildDiffTree(m1, m2), oldVal: v1}
		switch {
		case ok1 && ok2 && reflect.DeepEqual(v1, v2):
			node.status = sUnchanged
		case ok1 && ok2:
			node.status = sChanged
			node.newVal = v2
			node.oldSub = okv1
			node.newSub = okv2
		case ok1 && !ok2:
			node.status = sRemoved
			node.newVal = v2
		case !ok1 && ok2:
			node.status = sAdded
			node.newVal = v2
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func renderStylish(diffTree []*node) string {
	var sb strings.Builder
	sb.WriteString("{\n")

	renderStylishLevel(&sb, diffTree, 1, false)
	sb.WriteString("}")

	return sb.String()
}

func renderStylishLevel(sb *strings.Builder, nodes []*node, level int, marked bool) {

	plusChar, minusChar := " ", " "
	if !marked {
		plusChar = color.GreenString("+")
		minusChar = color.RedString("-")
	}

	for _, n := range nodes {
		switch n.status {
		case sUnchanged:
			addEntry(sb, " ", n.k, n.oldVal, n.childs, level)

		case sAdded:
			addEntry(sb, plusChar, n.k, n.newVal, n.childs, level)

		case sRemoved:
			addEntry(sb, minusChar, n.k, n.oldVal, n.childs, level)

		case sChanged:
			switch {
			case n.oldSub && n.newSub: // оба объекты
				fmt.Fprintf(sb, "%s%s: {\n", strings.Repeat("    ", level), n.k)
				renderStylishLevel(sb, n.childs, level+1, marked)
			case n.oldSub: // перестал быть объектом
				addEntry(sb, minusChar, n.k, nil, n.childs, level)
				addEntry(sb, plusChar, n.k, n.newVal, nil, level)
			case n.newSub: // стал объектом
				addEntry(sb, minusChar, n.k, n.oldVal, nil, level)
				addEntry(sb, plusChar, n.k, nil, n.childs, level)
			default: // оба не объекты
				addEntry(sb, minusChar, n.k, n.oldVal, nil, level)
				addEntry(sb, plusChar, n.k, n.newVal, nil, level)
			}
		}
	}
	if level > 1 {
		fmt.Fprintf(sb, "%s}\n", strings.Repeat("    ", level-1))
	}
}

func addEntry(sb *strings.Builder, sign, key string, value any, childs []*node, level int) {
	indent := strings.Repeat("    ", level-1) + "  "

	if len(childs) > 0 {
		fmt.Fprintf(sb, "%s%s %s: {\n", indent, sign, key)
		renderStylishLevel(sb, childs, level+1, true)
	} else {
		fmt.Fprintf(sb, "%s%s %s: %s\n", indent, sign, key, formatValue(value))
	}
}

func formatValue(v any) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}

func renderPlain(diffTree []*node) string {
	return ""
}
