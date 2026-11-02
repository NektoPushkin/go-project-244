package formatters

import (
	"encoding/json"
	"strings"
)

type patchNode map[string]any

func RenderJSONPatch(diffTree []*Node) (string, error) {
	jOut, err := json.MarshalIndent(buildJSONPatchStruct(diffTree, nil), "", "    ")
	return string(jOut), err
}

func buildJSONPatchStruct(nodes []*Node, path []string) []patchNode {
	patchElems := make([]patchNode, 0)
	for _, n := range nodes {
		switch n.Status {
		case StatusAdded:
			patchElems = append(patchElems, makeJSONPatchElem("add", makeJSONPatchPath(path, n.K), n.NewVal))
		case StatusRemoved:
			patchElems = append(patchElems, makeJSONPatchElem("remove", makeJSONPatchPath(path, n.K), nil))
		case StatusNested:
			if len(n.Childs) > 0 {
				patchElems = append(patchElems, buildJSONPatchStruct(n.Childs, append(path, n.K))...)
			}
		case StatusChanged:
			patchElems = append(patchElems, makeJSONPatchElem("replace", makeJSONPatchPath(path, n.K), n.NewVal))
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

var jsonPointerEscaper = strings.NewReplacer("~", "~0", "/", "~1")

func makeJSONPatchPath(path []string, curr string) string {
	tokens := append(append([]string{}, path...), curr)
	for i, tok := range tokens {
		tokens[i] = jsonPointerEscaper.Replace(tok)
	}
	return "/" + strings.Join(tokens, "/")
}
