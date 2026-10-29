package parsers

import (
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
			"setting2": float64(200),
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
			"abc": float64(12345),
			"deep": map[string]any{
				"id": float64(45),
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
			"setting2": 200,
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
			"abc": 12345,
			"deep": map[string]any{
				"id": 45,
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
