package formatters

import (
	"fmt"
	"strings"
)

func RenderPlain(diffTree []*Node) string {
	var sb strings.Builder
	renderPlainLevel(&sb, diffTree, make([]string, 0))
	return sb.String()
}

func renderPlainLevel(sb *strings.Builder, nodes []*Node, path []string) {
	curPath := strings.Join(path, ".")
	for i, n := range nodes {
		strPrefix := ""
		if i > 0 {
			strPrefix = "\n"
		}
		nodePath := n.K
		if curPath != "" {
			nodePath = curPath + "." + n.K
		}
		switch n.Status {
		case StatusAdded:
			fmt.Fprintf(sb, "%sProperty '%s' was added with value: %v", strPrefix, nodePath, formatPlainValue(n.NewVal))
		case StatusRemoved:
			fmt.Fprintf(sb, "%sProperty '%s' was removed", strPrefix, nodePath)
		case StatusChanged:
			switch {
			case n.OldSub && n.NewSub: // оба объекты
				renderPlainLevel(sb, n.Childs, append(path, n.K))
			case n.OldSub: // перестал быть объектом
				fmt.Fprintf(sb, "%sProperty '%s' was updated. From [complex value] to %s", strPrefix, nodePath, formatPlainValue(n.NewVal))
			case n.NewSub: // стал объектом
				fmt.Fprintf(sb, "%sProperty '%s' was updated. From %s to [complex value]", strPrefix, nodePath, formatPlainValue(n.OldVal))
			default: // оба не объекты
				fmt.Fprintf(sb, "%sProperty '%s' was updated. From %s to %s", strPrefix, nodePath, formatPlainValue(n.OldVal), formatPlainValue(n.NewVal))
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
