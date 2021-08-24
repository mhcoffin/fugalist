package fugalist

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"sort"
	"testing"
)

// These tests work as follows:
// Read an expression map (xml) to create doricolib.ScoreLib
// Pull out the expression map
// Convert some part of the expression map into a fugalist type (Axis, etc.)
// Test by
//   - converting the fugalist thing into a doricolib thing
//   - comparing to the original
//

func ReadDoricolib(t *testing.T, name string) *doricolib.ScoreLib {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("test_input/%s.doricolib", name))
	if err != nil {
		t.Fatalf("failed to read %s.doricolib", name)
	}
	scoreLib := &doricolib.ScoreLib{}
	err = xml.Unmarshal(bytes, scoreLib)
	if err != nil {
		t.Fatalf("failed to unmarshall %s.doricolib", name)
	}
	return scoreLib
}

func ReadProject(t *testing.T, name string) *Project {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("test_input/%s.project.json", name))
	if err != nil {
		t.Fatalf("failed to read %s.json", name)
	}
	project := &Project{}
	err = json.Unmarshal(bytes, project)
	if err != nil {
		t.Fatalf("failed to unmarshall project %s.project.json", name)
	}
	return project
}

func getXmap(scorelib *doricolib.ScoreLib) *doricolib.ExpressionMap {
	if len(scorelib.ExpressionMaps.Entities.Contents) == 0 {
		return nil
	}
	return &scorelib.ExpressionMaps.Entities.Contents[0]
}

func getPtMap(t *testing.T, xmap *doricolib.ExpressionMap) PtMap {
	ptMap, err := BuildPtMap(xmap)
	if err != nil {
		t.Fatalf("failed to build ptMap")
	}
	return ptMap
}

func getCombos(ptMap PtMap) []string {
	combos := make([]string, len(ptMap))
	k := 0
	for key := range ptMap {
		combos[k] = key
		k++
	}
	return combos
}

func getSortedAxes(axes map[AxisId]Axis) []Axis {
	result := make([]Axis, len(axes))
	k := 0
	for _, value := range axes {
		result[k] = value
		k++
	}
	sort.Slice(result, func(a int, b int) bool { return result[a].SortOrder < result[b].SortOrder })
	return result
}

func assertEqualAxes(t *testing.T, expected map[AxisId]Axis, actual map[AxisId]Axis) {
	expectedSorted := getSortedAxes(expected)
	actualSorted := getSortedAxes(actual)
	assert.Equal(t, len(expectedSorted), len(actualSorted))
	for k := 0; k < len(expectedSorted); k++ {
		assert.Equal(t, expectedSorted[k].Name, actualSorted[k].Name)
		assert.Equal(t, len(expectedSorted[k].Techniques), len(actualSorted[k].Techniques))
		for n := 0; n < len(expectedSorted[k].Techniques); n++ {
			assert.Equal(t, expectedSorted[k].Techniques[n].Name, actualSorted[k].Techniques[n].Name)
		}
	}
}

func TestFindAxes2(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"ParseTest1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score := ReadDoricolib(t, test.name)
			xmap := getXmap(score)
			ptMap := getPtMap(t, xmap)
			combos := getCombos(ptMap)
			axes := FindAxes(combos)

			project := ReadProject(t, test.name)
			assertEqualAxes(t, project.Axes, axes)
		})
	}
}

func vstSoundsByName(vstSoundsById map[VstSoundId]*VstSound) map[string]*VstSound {
	result := make(map[string]*VstSound)
	for _, vstSound := range vstSoundsById {
		result[vstSound.Name] = vstSound
	}
	return result
}

func assertEqualVstSounds(t *testing.T, expected map[VstSoundId]*VstSound, actual map[VstSoundId]*VstSound) {
	assert.Equal(t, len(expected), len(actual))
	assert.Equal(t, expected, actual)
}
