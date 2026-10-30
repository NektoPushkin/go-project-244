package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type dataParser func([]byte, any) error

var parsers = map[string]dataParser{
	"json": json.Unmarshal,
	"yml":  yaml.Unmarshal,
	"yaml": yaml.Unmarshal,
}

// ParseFile читает файлы конфигураций разбирает их
// в зависимости от расширения (json, yml, yaml).
func ParseFile(path string) (map[string]any, error) {
	realPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(realPath), "."))

	parse, ok := parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported config type: %s", path)
	}

	cont, err := os.ReadFile(realPath)
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if err := parse(cont, &data); err != nil {
		return nil, err
	}

	return data, nil
}
