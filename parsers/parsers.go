package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type dataParser func([]byte, any) error

var parsers = map[string]dataParser{
	"json": parseJSON,
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
		return nil, fmt.Errorf("error parsing '%s': %w", path, err)
	}

	if data != nil {
		data = normalizeNumbers(data).(map[string]any)
	}

	return data, nil
}

func parseJSON(b []byte, dst any) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	return decoder.Decode(dst)
}

func normalizeNumbers(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = normalizeNumbers(val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = normalizeNumbers(val)
		}
		return v
	case json.Number: //json
		if i, err := v.Int64(); err == nil {
			return i
		}
		if u, err := strconv.ParseUint(v.String(), 10, 64); err == nil {
			return u
		}
		f, _ := v.Float64()
		return f
	case int:
		return int64(v)
	case uint64:
		if v <= math.MaxInt64 {
			return int64(v)
		}
		return v
	default:
		return v
	}
}
