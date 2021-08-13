package fugalist

import (
	"encoding/xml"
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func readExpressionMap(filename string) *doricolib.ExpressionMap {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}
	scoreLib := &doricolib.ScoreLib{}
	err = xml.Unmarshal(bytes, scoreLib)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshall file (%s): %w", filename, err))
	}
	return &scoreLib.ExpressionMaps.Entities.Contents[0]
}

func TestCanonicalizeTechniqueString(t *testing.T) {
	tests := []struct {
		name     string
		orig     string
		expected string
	}{
		{"single", "pt.legato", "pt.legato"},
		{"sorted", "pt.legato+pt.nonVibrato", "pt.legato+pt.nonVibrato"},
		{"reversed", "pt.nonVibrato+pt.legato", "pt.legato+pt.nonVibrato"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, CanonicalizeTechniqueString(test.orig))
		})
	}
}

func TestFormatMidiEvents(t *testing.T) {
	tests := []struct {
		name     string
		actions  []doricolib.SwitchAction
		expected string
	}{
		{"empty", []doricolib.SwitchAction{}, ""},
		{"single KS", []doricolib.SwitchAction{
			{
				Type:   "kKeySwitch",
				Param1: "13",
			},
		}, "KS13"},
		{"single PC", []doricolib.SwitchAction{
			{
				Type:   "kProgramChange",
				Param1: "7",
			},
		}, "PC7"},
		{"single CC", []doricolib.SwitchAction{
			{
				Type:   "kControlChange",
				Param1: "3",
				Param2: "17",
			},
		}, "CC3=17"},
		{"multiple", []doricolib.SwitchAction{
			{
				Type:   "kKeySwitch",
				Param1: "13",
			},
			{
				Type:   "kProgramChange",
				Param1: "7",
			},
			{
				Type:   "kControlChange",
				Param1: "3",
				Param2: "17",
			},
		}, "KS13, PC7, CC3=17"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatMidiEvents(test.actions))
		})
	}
}

func TestFormatMidiDynamic(t *testing.T) {
	tests := []struct {
		name     string
		vtype    doricolib.VolumeType
		rng      string
		expected string
	}{
		{"velocity full range", doricolib.VolumeType{Type: "kNoteVelocity"}, "0,127", "velocity"},
		{"velocity part range", doricolib.VolumeType{Type: "kNoteVelocity"}, "10,110", "velocity 10:110"},
		{"cc full range", doricolib.VolumeType{Type: "kCC", Param1: "13"}, "0,127", "CC13"},
		{"cc part range", doricolib.VolumeType{Type: "kCC", Param1: "13"}, "10,30", "CC13 10:30"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatMidiDynamic(test.vtype, test.rng))
		})
	}
}

func TestFormatLengthFactor(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected float64
	}{
		{"empty", "", 100},
		{"one", "1.0", 100},
		{"fraction", "0.85", 85},
		{"zero", "0.0", 0},
		{"larger", "1.05", 105},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatLengthFactor(test.value))
		})
	}
}

func TestFormatTranspose(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected float64
	}{
		{"empty", 0, 0.0},
		{"empty", 1, 1.0},
		{"empty", -1, -1.0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatTranspose(test.value))

		})
	}
}

func TestFormatBranch(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"empty", "", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, FormatBranch(test.value))
		})
	}
}

func TestBuildPtMap(t *testing.T) {
	tests := []struct {
		name     string
		expected PtMap
	}{
		{
			name: "Ref",
			expected: PtMap{
				"pt.legato": {
					"NoteLength <= medium": {
						On:    "KS25, PC6, CC1=64",
						Dyn:   "velocity 1:127",
						Len:   100,
						Trans: 0,
					},
					"NoteLength > medium": {
						On:    "KS26, PC6, CC1=64",
						Dyn:   "CC2 1:120",
						Len:   95,
						Trans: -1,
					},
				},
				"pt.marcato+pt.nonVibrato+pt.plucked": {
					"": {
						On:    "KS24, PC13, CC7=23",
						Dyn:   "velocity 1:127",
						Len:   100,
						Trans: 0,
					},
				},
				"pt.natural": {
					"NoteLength < medium": {
						On:    "KS12=120, KS24, PC15, CC4=64",
						Dyn:   "velocity 10:120",
						Len:   100,
						Trans: 0,
					},
					"NoteLength >= long": {
						On:    "KS12=120, KS24, PC13, CC4=64",
						Dyn:   "CC2 10:120",
						Len:   100,
						Trans: 0,
					},
					"NoteLength >= medium AND NoteLength < long": {
						On:    "KS12=120, KS24, PC13, CC4=64",
						Dyn:   "velocity 10:120",
						Trans: 0,
						Len:   100,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			xmap := readExpressionMap(fmt.Sprintf("test_input/%s.doricolib", test.name))
			ptMap, err := BuildPtMap(*xmap)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, ptMap)
		})
	}
}
