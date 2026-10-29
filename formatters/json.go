package formatters

import (
	"encoding/json"
	"strings"
)

type patchNode map[string]any

func RenderJSON(diffTree []*Node) string {
	jOut, err := json.MarshalIndent(buildJSONPatchStruct(diffTree, nil), "", "    ")
	if err != nil {
		return "" //гасим ошибку для совместимости
	}
	return string(jOut)
}

func buildJSONPatchStruct(nodes []*Node, path []string) []patchNode {
	patchElems := make([]patchNode, 0)
	for _, n := range nodes {
		switch n.Status {
		case StatusAdded:
			patchElems = append(patchElems, makeJSONPatchElem("add", makeJSONPatchPath(path, n.K), n.NewVal))
		case StatusRemoved:
			patchElems = append(patchElems, makeJSONPatchElem("remove", makeJSONPatchPath(path, n.K), nil))
		case StatusChanged:
			if n.OldSub && n.NewSub {
				patchElems = append(patchElems, buildJSONPatchStruct(n.Childs, append(path, n.K))...)
			} else {
				patchElems = append(patchElems, makeJSONPatchElem("replace", makeJSONPatchPath(path, n.K), n.NewVal))
			}
		case StatusUnchanged:
			continue //чтобы не потерять
		}
	}
	return patchElems
}

func makeJSONPatchElem(action, path string, val any) patchNode {
	pn := patchNode{"op": action, "path": path}
	if action != "remove" {
		pn["value"] = val
	}
	return pn
}

func makeJSONPatchPath(path []string, curr string) string {
	return "/" + strings.Join(append(append([]string{}, path...), curr), "/")
}
