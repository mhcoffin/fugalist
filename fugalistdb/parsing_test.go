package fugalistdb

import (
	"github.com/mhcoffin/go-doricolib/doricolib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMidiParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"basic", "CC1=2", []string{"1", "2"}},
		{"ws", " CC 2 = 17 ", []string{"2", "17"}},
		{"lower case", " cc2 = 17 ", []string{"2", "17"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := CcPat.FindStringSubmatch(test.input)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}
func TestKsParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"basic", "KS1=2", []string{"1", "2"}},
		{"ws", " KS 2 = 17 ", []string{"2", "17"}},
		{"lower case", " ks 2", []string{"2", ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := KsPat.FindStringSubmatch(test.input)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}

func TestNoteParser(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
	}{
		{"C0", []string{"C", "", "0"}},
		{" C2", []string{"C", "", "2"}},
		{"C#0", []string{"C", "#", "0"}},
		{"d#3", []string{"d", "#", "3"}},
		{"d#-1", []string{"d", "#", "-1"}},
		{"db -1", []string{"d", "b", "-1"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := NotePat.FindStringSubmatch(test.name)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}

func TestPcParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"basic", "PC1", []string{"1"}},
		{"ws", " pc 2", []string{"2"}},
		{"lower case", " pc 2", []string{"2"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := PcPat.FindStringSubmatch(test.input)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}

func TestCcPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"CC without range", "CC7", []string{"7", "", ""}},
		{"CC with spaces", "  CC 7   ", []string{"7", "", ""}},
		{"CC with spaces", "  CC 77   ", []string{"77", "", ""}},
		{"CC with range", "CC3 0-127", []string{"3", "0", "127"}},
		{"CC with range", "CC3 0-127", []string{"3", "0", "127"}},
		{"CC with range and space", " CC 3 0 - 127 ", []string{"3", "0", "127"}},
		{"lower case", " cc 3 0 - 127 ", []string{"3", "0", "127"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := CcPattern.FindStringSubmatch(test.input)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}

func TestVelPattern(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Vel without range", "VEL", []string{"", ""}},
		{"Vel lower case", "vel", []string{"", ""}},
		{"Vel with range", "Vel 10-90", []string{"10", "90"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := VelPattern.FindStringSubmatch(test.input)
			assert.Equal(t, test.expected, actual[1:])
		})
	}
}

func TestParseVolumeSpec(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedVolType *doricolib.VolumeType
		expectedRange   string
	}{
		{"CC", "cc3", &doricolib.VolumeType{Type: "kCC", Param1: "3"}, "0,127"},
		{"CC", "cc3 10-100", &doricolib.VolumeType{Type: "kCC", Param1: "3"}, "10,100"},
		{"Vel", "vel", &doricolib.VolumeType{Type: "kNoteVelocity", Param1: "0"}, "0,127"},
		{"Vel", "vel 11-33", &doricolib.VolumeType{Type: "kNoteVelocity", Param1: "0"}, "11,33"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vt, rng, err := ParseVolumeSpec(test.input)
			assert.Nil(t, err)
			assert.Equal(t, test.expectedVolType, vt)
			assert.Equal(t, test.expectedRange, rng)
		})
	}
}

func TestNote(t *testing.T) {
	tests := []struct {
		name          string
		noteName      string
		incidental    string
		octave        string
		middleCOctave int
		expected      int
	}{
		{"middle C", "C", "", "4", 4, 60},
		{"middle C#", "C", "#", "4", 4, 61},
		{"middle Db", "D", "b", "4", 4, 61},
		{"middle D", "D", "", "4", 4, 62},
		{"middle D#", "D", "#", "4", 4, 63},
		{"middle Eb", "E", "b", "4", 4, 63},
		{"middle E", "E", "", "4", 4, 64},
		{"middle F", "F", "", "4", 4, 65},
		{"middle F#", "F", "#", "4", 4, 66},
		{"middle Gb", "G", "b", "4", 4, 66},
		{"middle G", "G", "", "4", 4, 67},
		{"middle G#", "G", "#", "4", 4, 68},
		{"middle Ab", "A", "b", "4", 4, 68},
		{"middle A", "A", "", "4", 4, 69},
		{"middle A#", "A", "#", "4", 4, 70},
		{"middle Bb", "B", "b", "4", 4, 70},
		{"middle B", "B", "", "4", 4, 71},

		{"middle C, middle C octave = 3", "C", "", "3", 3, 60},
		{"middle C#", "C", "#", "3", 3, 61},

		{"zero octave", "C", "#", "0", 4, 13},
		{"-1 octave", "C", "#", "-1", 4, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n, err := note(test.noteName, test.incidental, test.octave, test.middleCOctave)
			assert.Nil(t, err)
			assert.Equal(t, n, test.expected)
		})
	}
}
