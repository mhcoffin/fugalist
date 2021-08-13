package fugalist

import (
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"math"
	"sort"
	"strconv"
	"strings"
)

type PlayData struct {
	On    string
	Off   string
	Dyn   string
	Len   float64
	Trans float64
}

type BrMap map[string]PlayData

type PtMap map[string]BrMap

func BuildPtMap(xmap doricolib.ExpressionMap) (PtMap, error) {
	result := make(PtMap)
	for _, combo := range xmap.Combinations.Combos {
		tids := CanonicalizeTechniqueString(combo.TechniqueIDs)
		cond := FormatBranch(combo.ConditionString)
		playData := PlayData{
			On:    FormatMidiEvents(combo.SwitchOnActions.SwitchOnActions),
			Off:   FormatMidiEvents(combo.SwitchOffActions.SwitchOffActions),
			Dyn:   FormatMidiDynamic(combo.VolumeType, combo.VelocityRange),
			Len:   FormatLengthFactor(combo.LengthFactor),
			Trans: FormatTranspose(combo.Transpose),
		}
		branchMap, exists := result[tids]
		if exists {
			branchMap[cond] = playData
		} else {
			result[tids] = make(BrMap)
			result[tids][cond] = playData
		}
	}
	return result, nil
}

var branchReplacer = strings.NewReplacer(
	"kVeryShort", "veryShort",
	"kShort", "short",
	"kMedium", "medium",
	"kLong", "long",
	"kVeryLong", "veryLong",
)

func FormatBranch(br string) string {
	return branchReplacer.Replace(br)
}

func FormatTranspose(transpose int) float64 {
	return float64(transpose)
}

func FormatLengthFactor(factor string) float64 {
	fl, err := strconv.ParseFloat(factor, 32)
	if err != nil {
		return 100.0
	}
	return math.Round(fl * 100.0)
}

func FormatMidiDynamic(volumeType doricolib.VolumeType, vrange string) string {
	rng := ""
	parsedRange := strings.Split(vrange, ",")
	if len(parsedRange) == 2 && (parsedRange[0] != "0" || parsedRange[1] != "127") {
		rng = fmt.Sprintf(" %s:%s", parsedRange[0], parsedRange[1])
	}
	if volumeType.Type == "kNoteVelocity" {
		return fmt.Sprintf("velocity%s", rng)
	} else if volumeType.Type == "kCC" {
		return fmt.Sprintf("CC%s%s", volumeType.Param1, rng)
	} else {
		return fmt.Sprintf("velocity")
	}
}

func FormatMidiEvents(actions []doricolib.SwitchAction) string {
	var result []string
	for _, action := range actions {
		switch action.Type {
		case "kKeySwitch":
			// TODO: convert to notes?
			if action.Param2 != "" && action.Param2 != "127" {
				result = append(result, fmt.Sprintf("KS%s=%s", action.Param1, action.Param2))
			} else {
				result = append(result, fmt.Sprintf("KS%s", action.Param1))
			}
		case "kProgramChange":
			result = append(result, fmt.Sprintf("PC%s", action.Param1))
		case "kControlChange":
			result = append(result, fmt.Sprintf("CC%s=%s", action.Param1, action.Param2))
		}
	}
	return strings.Join(result, ", ")
}

func CanonicalizeTechniqueString(tids string) string {
	t := strings.Split(tids, "+")
	sort.Strings(t)
	return strings.Join(t, "+")
}
