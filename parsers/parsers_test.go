package parsers

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fixtureFlatJSON1  = "../testdata/fixture/json.flat.1.json"
	fixtureFlatYAML1  = "../testdata/fixture/yaml.flat.1.yml"
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

func TestParseFile_FlatJSON(t *testing.T) {
	t.Parallel()
	color.NoColor = true
	want := map[string]any{
		"host":    "hexlet.io",
		"timeout": float64(50),
		"proxy":   "123.234.53.22",
		"follow":  false,
	}

	got, err := ParseFile(fixtureFlatJSON1)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseFile_FlatYAML(t *testing.T) {
	t.Parallel()
	color.NoColor = true
	want := map[string]any{
		"host":    "hexlet.io",
		"timeout": 50,
		"proxy":   "123.234.53.22",
		"follow":  false,
	}

	got, err := ParseFile(fixtureFlatYAML1)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestParseFileMissingFile(t *testing.T) {
	t.Parallel()
	_, err := ParseFile("testdata/fixture/file-not-exist.json")
	assert.Error(t, err)
}
