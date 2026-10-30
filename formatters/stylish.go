package formatters

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	plusChar  = color.GreenString("+")
	minusChar = color.RedString("-")
)

func RenderStylish(diffTree []*Node) string {
	var sb strings.Builder
	sb.WriteString("{\n")

	renderStylishLevel(&sb, diffTree, 1)
	sb.WriteString("}")

	return sb.String()
}

func renderStylishLevel(sb *strings.Builder, nodes []*Node, level int) {
	for _, n := range nodes {
		switch n.Status {
		case StatusAdded:
			addStylishEntry(sb, plusChar, n.K, n.NewVal, n.Childs, level)

		case StatusRemoved:
			addStylishEntry(sb, minusChar, n.K, n.OldVal, n.Childs, level)

		case StatusChanged:
			addStylishEntry(sb, minusChar, n.K, n.OldVal, nil, level)
			addStylishEntry(sb, plusChar, n.K, n.NewVal, nil, level)

		default:
			addStylishEntry(sb, " ", n.K, n.OldVal, n.Childs, level)
		}

	}
	if level > 1 {
		fmt.Fprintf(sb, "%s}\n", strings.Repeat("    ", level-1))
	}
}

func addStylishEntry(sb *strings.Builder, sign, key string, value any, childs []*Node, level int) {
	indent := strings.Repeat("    ", level-1) + "  "

	if len(childs) > 0 {
		fmt.Fprintf(sb, "%s%s %s: {\n", indent, sign, key)
		renderStylishLevel(sb, childs, level+1)
	} else {
		fmt.Fprintf(sb, "%s%s %s: %s\n", indent, sign, key, formatStylishValue(value, level+1))
	}
}

func formatStylishValue(v any, level int) string {
	switch v := v.(type) {
	case nil:
		return "null"
	case map[string]any:
		var sb strings.Builder
		sb.WriteString("{\n")
		for k, v2 := range v {
			fmt.Fprintf(&sb, "%s%s: %s\n", strings.Repeat("    ", level), k, formatStylishValue(v2, level+1))
		}
		sb.WriteString(strings.Repeat("    ", level-1))
		sb.WriteString("}")
		return sb.String()
	case float64:
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
