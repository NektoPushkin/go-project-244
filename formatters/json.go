package formatters

import (
	"encoding/json"
)

func RenderJSON(diffTree []*Node) (string, error) {

	data, err := json.Marshal(Node{K: "", Status: "root", Childs: diffTree})
	return string(data), err
}
