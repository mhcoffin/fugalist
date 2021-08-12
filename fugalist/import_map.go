package fugalist

import (
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"sort"
	"strconv"
	"strings"
)

type PlayData struct {
	On string
	Off string
	Dyn string
	Len float64
	Trans float64
}

type BrMap map[string]PlayData

type PtMap map[string]BrMap

func BuildPtMap(xmap doricolib.ExpressionMap) (PtMap, error) {
	result := make(PtMap)
	for _, combo := range xmap.Combinations.Combos {
		tids := CanonicalizeTechniqueString(combo.TechniqueIDs)
		cond := combo.ConditionString
		brmap, exists := result[tids]
		playData := PlayData {
			On: FormatMidiEvents(combo.SwitchOnActions.SwitchOnActions),
			Off: FormatMidiEvents(combo.SwitchOffActions.SwitchOffActions),
			Dyn: FormatMidiDynamic(combo.VolumeType),
			Len: FormatLengthFactor(combo.LengthFactor),
			Trans: FormatTranspose(combo.Transpose),
		}
		if exists {
			brmap[cond] = PlayData{}
		} else {
			result[tids] = make(BrMap)
			result[tids][cond] = PlayData{}
		}
	}
	return result, nil
}

func FormatTranspose(transpose int) float64 {
	return float64(transpose)
}

func FormatLengthFactor(factor string) float64 {
	fl, err := strconv.ParseFloat(factor, 32)
	if err != nil {
		return 100.0
	}
	return fl * 100.0
}

func FormatMidiDynamic(volumeType doricolib.VolumeType) string {
	if volumeType.Type == "kNoteVelocity" {
		return fmt.Sprintf("velocity")
	}
}

func FormatMidiEvents(actions []doricolib.SwitchAction) string {

}

func CanonicalizeTechniqueString(tids string) string {
	t := strings.Split(tids, "+")
	sort.Strings(t)
	return strings.Join(t, "+")
}