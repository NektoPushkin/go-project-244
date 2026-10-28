package code

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fixtureFlatJSON1  = "testdata/fixture/json.flat.1.json"
	fixtureFlatJSON2  = "testdata/fixture/json.flat.2.json"
	fixtureBrokenJSON = "testdata/fixture/json.broken.json"
)

func TestGenDiffStylish_Flat(t *testing.T) {
	t.Parallel()
	color.NoColor = true
	const want = `{
  - follow: false
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`

	got, err := GenDiff(fixtureFlatJSON1, fixtureFlatJSON2, "stylish")
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestGenDiffIdenticalFiles(t *testing.T) {
	t.Parallel()
	got, err := GenDiff(fixtureFlatJSON1, fixtureFlatJSON1, "stylish")
	require.NoError(t, err)
	assert.Equal(t, noDifferenceFound, got)
}

func TestGenDiffBrokenInput(t *testing.T) {
	t.Parallel()
	_, err := GenDiff(fixtureBrokenJSON, fixtureFlatJSON2, "stylish")
	assert.Error(t, err)
}

func TestGenDiffUnknownFormat(t *testing.T) {
	t.Parallel()
	got, err := GenDiff(fixtureFlatJSON1, fixtureFlatJSON2, "raw")
	require.ErrorIs(t, err, ErrUnsupportedFormat)
	assert.Empty(t, got)
}

func TestRenderStylish_Flat(t *testing.T) {
	t.Parallel()
	color.NoColor = true
	diff := diff{
		"added":   {nil, true},
		"removed": {1, nil},
		"changed": {"was", "now"},
		"same":    {2, 2},
	}
	const want = `{
  + added: true
  - changed: was
  + changed: now
  - removed: 1
}`

	assert.Equal(t, want, renderStylish(diff))
}
