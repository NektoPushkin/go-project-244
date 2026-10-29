package formatters

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func RenderStylish(diffTree []*Node) string {
	var sb strings.Builder
	sb.WriteString("{\n")

	renderStylishLevel(&sb, diffTree, 1, false)
	sb.WriteString("}")

	return sb.String()
}

func renderStylishLevel(sb *strings.Builder, nodes []*Node, level int, marked bool) {

	plusChar, minusChar := " ", " "
	if !marked {
		plusChar = color.GreenString("+")
		minusChar = color.RedString("-")
	}

	for _, n := range nodes {
		switch n.Status {
		case StatusUnchanged:
			addStylishEntry(sb, " ", n.K, n.OldVal, n.Childs, level)

		case StatusAdded:
			addStylishEntry(sb, plusChar, n.K, n.NewVal, n.Childs, level)

		case StatusRemoved:
			addStylishEntry(sb, minusChar, n.K, n.OldVal, n.Childs, level)

		case StatusChanged:
			switch {
			case n.OldSub && n.NewSub: // оба объекты
				fmt.Fprintf(sb, "%s%s: {\n", strings.Repeat("    ", level), n.K)
				renderStylishLevel(sb, n.Childs, level+1, marked)
			case n.OldSub: // перестал быть объектом
				addStylishEntry(sb, minusChar, n.K, nil, n.Childs, level)
				addStylishEntry(sb, plusChar, n.K, n.NewVal, nil, level)
			case n.NewSub: // стал объектом
				addStylishEntry(sb, minusChar, n.K, n.OldVal, nil, level)
				addStylishEntry(sb, plusChar, n.K, nil, n.Childs, level)
			default: // оба не объекты
				addStylishEntry(sb, minusChar, n.K, n.OldVal, nil, level)
				addStylishEntry(sb, plusChar, n.K, n.NewVal, nil, level)
			}
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
		renderStylishLevel(sb, childs, level+1, true)
	} else {
		fmt.Fprintf(sb, "%s%s %s: %s\n", indent, sign, key, formatStylishValue(value))
	}
}

func formatStylishValue(v any) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v)
}
