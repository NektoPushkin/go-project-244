package code

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"code/formatters"
)

const (
	fixtureJSON1      = "testdata/fixture/json.1.json"
	fixtureJSON2      = "testdata/fixture/json.2.json"
	fixtureBrokenJSON = "testdata/fixture/json.broken.json"
	fixtureYAML1      = "testdata/fixture/yaml.1.yaml"
	fixtureYAML2      = "testdata/fixture/yaml.2.yml"
	fixtureBrokenYAML = "testdata/fixture/yaml.broken.yml"

	wantOutput = `{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
      + setting4: blah blah
      + setting5: {
            key5: value5
        }
        setting6: {
            doge: {
              - wow: 
              + wow: so much
            }
            key: value
          + ops: vops
        }
    }
    group1: {
      - baz: bas
      + baz: bars
        foo: bar
      - nest: {
            key: value
        }
      + nest: str
    }
  - group2: {
        abc: 12345
        deep: {
            id: 45
        }
    }
  + group3: {
        deep: {
            id: {
                number: 45
            }
        }
        fee: 100500
    }
}`
)

func TestGenDiffStylish_JSON(t *testing.T) {
	t.Parallel()
	color.NoColor = true

	got, err := GenDiff(fixtureJSON1, fixtureJSON2, "stylish")
	require.NoError(t, err)
	assert.Equal(t, wantOutput, got)
}

func TestGenDiffStylish_YAML(t *testing.T) {
	t.Parallel()
	color.NoColor = true

	got, err := GenDiff(fixtureYAML1, fixtureYAML2, "stylish")
	require.NoError(t, err)
	assert.Equal(t, wantOutput, got)
}

func TestGenDiffIdenticalFiles(t *testing.T) {
	t.Parallel()
	got, err := GenDiff(fixtureJSON1, fixtureJSON1, "stylish")
	require.NoError(t, err)
	assert.Equal(t, noDifferenceFound, got)
}

func TestGenDiffBrokenInput(t *testing.T) {
	t.Parallel()
	_, err := GenDiff(fixtureBrokenJSON, fixtureJSON2, "stylish")
	assert.Error(t, err)
}

func TestGenDiffUnknownFormat(t *testing.T) {
	t.Parallel()
	got, err := GenDiff(fixtureJSON1, fixtureJSON2, "raw")
	require.ErrorIs(t, err, formatters.ErrUnsupportedFormat)
	assert.Empty(t, got)
}
