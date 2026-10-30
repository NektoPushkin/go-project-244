package code

import (
	"reflect"
	"slices"

	"code/formatters"
	"code/parsers"
)

var (
	noDifferenceFound = "no differences found"
)

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
	return formatters.Render(diff, format)
}

func buildDiffTree(config1, config2 map[string]any) []*formatters.Node {
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
	nodes := make([]*formatters.Node, 0, len(keys))

	for _, k := range keys {
		v1, ok1 := config1[k]
		v2, ok2 := config2[k]

		m1, isMap1 := v1.(map[string]any)
		m2, isMap2 := v2.(map[string]any)

		switch {
		case ok1 && ok2:
			switch {
			case isMap1 && isMap2:
				nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusNested, Childs: buildDiffTree(m1, m2)})
			case isMap1 || isMap2:
				nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusChanged, OldVal: v1, NewVal: v2, NewSub: isMap2, OldSub: isMap1})
			default:
				if v1 != v2 {
					nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusChanged, OldVal: v1, NewVal: v2})
				} else {
					nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusUnchanged, OldVal: v1})
				}
			}
		case ok1 && !ok2:
			nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusRemoved, OldVal: v1})
		case !ok1 && ok2:
			nodes = append(nodes, &formatters.Node{K: k, Status: formatters.StatusAdded, NewVal: v2})
		}
	}

	return nodes
}
