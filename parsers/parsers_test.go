package parsers

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fixtureJSON1      = "../testdata/fixture/json.1.json"
	fixtureYAML1      = "../testdata/fixture/yaml.1.yaml"
	fixtureBrokenJSON = "../testdata/fixture/json.broken.json"
	fixtureBrokenYAML = "../testdata/fixture/yaml.broken.yml"
)

func TestParseFileBrokenJSON(t *testing.T) {
	t.Parallel()
	_, err := ParseFile(fixtureBrokenJSON)
	require.Error(t, err)
}

func TestParseFileBrokenYAML(t *testing.T) {
	t.Parallel()
	_, err := ParseFile(fixtureBrokenYAML)
	require.Error(t, err)
}

func TestParseFile_JSON(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": int64(200),
			"setting3": true,
			"setting6": map[string]any{
				"key": "value",
				"doge": map[string]any{
					"wow": "",
				},
			},
		},
		"group1": map[string]any{
			"baz": "bas",
			"foo": "bar",
			"nest": map[string]any{
				"key": "value",
			},
		},
		"group2": map[string]any{
			"abc": int64(12345),
			"deep": map[string]any{
				"id": int64(45),
			},
		},
	}

	got, err := ParseFile(fixtureJSON1)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseFile_YAML(t *testing.T) {
	t.Parallel()
	want := map[string]any{
		"common": map[string]any{
			"setting1": "Value 1",
			"setting2": int64(200),
			"setting3": true,
			"setting6": map[string]any{
				"key": "value",
				"doge": map[string]any{
					"wow": "",
				},
			},
		},
		"group1": map[string]any{
			"baz": "bas",
			"foo": "bar",
			"nest": map[string]any{
				"key": "value",
			},
		},
		"group2": map[string]any{
			"abc": int64(12345),
			"deep": map[string]any{
				"id": int64(45),
			},
		},
	}

	got, err := ParseFile(fixtureYAML1)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseFileMissingFile(t *testing.T) {
	t.Parallel()
	_, err := ParseFile("../testdata/fixture/file-not-exist.json")
	assert.Error(t, err)
}

func TestNormalizeNumbers(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   any
		want any
	}{
		{"jsonNumber_int", json.Number("42"), int64(42)},
		{"jsonNumber_float", json.Number("1.5"), 1.5},
		{"jsonNumber_uint64", json.Number("18446744073709551615"), uint64(math.MaxUint64)},
		{"int", 42, int64(42)},
		{"uint64_small", uint64(42), int64(42)},
		{"uint64_max", uint64(math.MaxUint64), uint64(math.MaxUint64)},
		{"float_pass", 1.5, 1.5},
		{"float_huge", json.Number("2e400"), "2e400"},
		{"string_pass", "text", "text"},
		{"bool_pass", true, true},
		{"nil_pass", nil, nil},
		{
			"nested",
			map[string]any{"a": []any{uint64(7), map[string]any{"b": json.Number("1.25")}}},
			map[string]any{"a": []any{int64(7), map[string]any{"b": 1.25}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, normalizeNumbers(tc.in))
		})
	}
}
