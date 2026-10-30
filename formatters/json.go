package formatters

import (
	"encoding/json"
)

func RenderJSON(diffTree []*Node) string {

	data, err := json.Marshal(Node{K: "", Status: "root", Childs: diffTree})
	if err != nil {
		return "" //гасим ошибку для совместимости
	}
	return string(data)
}
