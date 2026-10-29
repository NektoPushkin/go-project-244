package formatters

import (
	"fmt"
	"strings"
)

func RenderPlain(diffTree []*Node) string {
	var sb strings.Builder
	renderPlainLevel(&sb, diffTree, make([]string, 0))
	return strings.TrimSpace(sb.String())
}

func renderPlainLevel(sb *strings.Builder, nodes []*Node, path []string) {
	curPath := strings.Join(path, ".")
	for _, n := range nodes {
		nodePath := n.K
		if curPath != "" {
			nodePath = curPath + "." + n.K
		}
		switch n.Status {
		case StatusAdded:
			fmt.Fprintf(sb, "Property '%s' was added with value: %v\n", nodePath, formatPlainValue(n.NewVal))
		case StatusRemoved:
			fmt.Fprintf(sb, "Property '%s' was removed\n", nodePath)
		case StatusChanged:
			switch {
			case n.OldSub && n.NewSub: // оба объекты
				renderPlainLevel(sb, n.Childs, append(path, n.K))
			case n.OldSub: // перестал быть объектом
				fmt.Fprintf(sb, "Property '%s' was updated. From [complex value] to %s\n", nodePath, formatPlainValue(n.NewVal))
			case n.NewSub: // стал объектом
				fmt.Fprintf(sb, "Property '%s' was updated. From %s to [complex value]\n", nodePath, formatPlainValue(n.OldVal))
			default: // оба не объекты
				fmt.Fprintf(sb, "Property '%s' was updated. From %s to %s\n", nodePath, formatPlainValue(n.OldVal), formatPlainValue(n.NewVal))
			}
		}
	}
}

func formatPlainValue(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case map[string]any:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
