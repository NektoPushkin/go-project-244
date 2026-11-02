package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/fatih/color"
)

var (
	plusChar  = color.GreenString("+")
	minusChar = color.RedString("-")
)

func RenderStylish(diffTree []*Node) (string, error) {
	var sb strings.Builder
	sb.WriteString("{\n")

	renderStylishLevel(&sb, diffTree, 1)
	sb.WriteString("}")

	return sb.String(), nil
}

func renderStylishLevel(sb *strings.Builder, nodes []*Node, level int) {
	indent := strings.Repeat("    ", level-1) + "  "
	for _, n := range nodes {
		switch n.Status {
		case StatusUnchanged:
			fmt.Fprintf(sb, "%s  %s: %s\n", indent, n.K, formatStylishValue(n.OldVal, level+1))

		case StatusNested:
			fmt.Fprintf(sb, "%s  %s: {\n", indent, n.K)
			renderStylishLevel(sb, n.Childs, level+1)

		case StatusAdded:
			fmt.Fprintf(sb, "%s%s %s: %s\n", indent, plusChar, n.K, formatStylishValue(n.NewVal, level+1))

		case StatusRemoved:
			fmt.Fprintf(sb, "%s%s %s: %s\n", indent, minusChar, n.K, formatStylishValue(n.OldVal, level+1))

		case StatusChanged:
			fmt.Fprintf(sb, "%s%s %s: %s\n", indent, minusChar, n.K, formatStylishValue(n.OldVal, level+1))
			fmt.Fprintf(sb, "%s%s %s: %s\n", indent, plusChar, n.K, formatStylishValue(n.NewVal, level+1))
		}

	}
	if level > 1 {
		fmt.Fprintf(sb, "%s}\n", strings.Repeat("    ", level-1))
	}
}

func formatStylishValue(v any, level int) string {
	switch v := v.(type) {
	case nil:
		return "null"
	case map[string]any:
		var sb strings.Builder
		sb.WriteString("{\n")
		for _, k := range slices.Sorted(maps.Keys(v)) {
			fmt.Fprintf(&sb, "%s%s: %s\n", strings.Repeat("    ", level), k, formatStylishValue(v[k], level+1))
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
