package formatters

import (
	"os"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	color.NoColor = true
	os.Exit(m.Run())
}

// sampleTree возвращает дерево, покрывающее все статусы:
// added, removed, changed (скаляр), unchanged и вложенные объекты.
func sampleTree() []*Node {
	return []*Node{
		{
			K:      "common",
			Status: StatusChanged,
			OldSub: true,
			NewSub: true,
			Childs: []*Node{
				{K: "follow", Status: StatusAdded, NewVal: false},
				{K: "setting1", Status: StatusUnchanged, OldVal: "Value 1"},
				{K: "setting2", Status: StatusRemoved, OldVal: 200},
				{K: "setting3", Status: StatusChanged, OldVal: true, NewVal: nil},
			},
		},
		{
			K:      "group2",
			Status: StatusRemoved,
			Childs: []*Node{
				{K: "abc", Status: StatusRemoved, OldVal: 12345},
			},
		},
		{
			K:      "group3",
			Status: StatusAdded,
			NewVal: map[string]any{"fee": 100500},
			Childs: []*Node{
				{K: "fee", Status: StatusAdded, NewVal: 100500},
			},
		},
	}
}

const wantStylish = `{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
      - setting3: true
      + setting3: null
    }
  - group2: {
        abc: 12345
    }
  + group3: {
        fee: 100500
    }
}`

const wantPlain = `Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
Property 'common.setting3' was updated. From true to null
Property 'group2' was removed
Property 'group3' was added with value: [complex value]
`

const wantJSON = `[
    {
        "op": "add",
        "path": "/common/follow",
        "value": false
    },
    {
        "op": "remove",
        "path": "/common/setting2"
    },
    {
        "op": "replace",
        "path": "/common/setting3",
        "value": null
    },
    {
        "op": "remove",
        "path": "/group2"
    },
    {
        "op": "add",
        "path": "/group3",
        "value": {
            "fee": 100500
        }
    }
]`

func TestRender_Stylish(t *testing.T) {
	got, err := Render(sampleTree(), "stylish")
	require.NoError(t, err)
	assert.Equal(t, wantStylish, got)
}

func TestRender_Plain(t *testing.T) {
	got, err := Render(sampleTree(), "plain")
	require.NoError(t, err)
	assert.Equal(t, wantPlain, got)

}
func TestRender_JSON(t *testing.T) {
	got, err := Render(sampleTree(), "json")
	require.NoError(t, err)
	assert.Equal(t, wantJSON, got)
}

func TestRender_UnknownFormat(t *testing.T) {
	got, err := Render(sampleTree(), "raw")
	assert.ErrorIs(t, err, ErrUnsupportedFormat)
	assert.Empty(t, got)
}

func TestRenderJSON(t *testing.T) {
	got := RenderJSON(sampleTree())
	assert.Equal(t, wantJSON, got)
}
