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
	Len   string
	Trans string
}

type BrMap map[string]PlayData

type PtMap map[string]BrMap

func BuildPtMap(xmap *doricolib.ExpressionMap) (PtMap, error) {
	result := make(PtMap)
	for _, combo := range xmap.Combinations.Combos {
		tids := CanonicalizeTechniqueString(combo.TechniqueIDs)
		cond := FormatBranch(combo.ConditionString)
		playData := PlayData{
			On:    FormatMidiEvents(combo.SwitchOnActions.SwitchOnActions),
			Off:   FormatMidiEvents(combo.SwitchOffActions.SwitchOffActions),
			Dyn:   FormatMidiDynamic(combo.VolumeType, combo.VelocityRange),
			Len:   FormatLengthFactor(combo.LengthFactor, combo.Flags),
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

func FormatTranspose(transpose int) string {
	return fmt.Sprintf("%d", transpose)
}

func FormatLengthFactor(factor string, flag int) string {
	if flag == 0 {
		return ""
	}
	fl, err := strconv.ParseFloat(factor, 32)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d", int(math.Round(fl*100.0)))
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

var usualPts = map[string]bool{
	"pt.natural":       true,
	"pt.staccato":      true,
	"pt.staccatissimo": true,
	"pt.tenuto":        true,
	"pt.portato":       true,
	"pt.legato":        true,
	"pt.marcato":       true,
	"pt.nonVibrato":    true,
}

func FindExtraTechniques(combos []string) []string {
	extras := []string{}
	for _, combo := range combos {
		pts := strings.Split(combo, "+")
		for _, pt := range pts {
			if !usualPts[pt] {
				extras = append(extras, pt)
			}
		}
	}
	return extras
}

func FugalistTechnique(id string) Technique {
	if id == "pt.normal" {
		return Technique{
			Id:   Uniq(),
			Name: "Normal",
		}
	}
	doricoTechnique := doricolib.GetTechniqueById(id)
	return Technique{Id: Uniq(), Name: doricoTechnique.Name}
}

// DefaultAxes returns a set of default axes with new unique IDs.
func DefaultAxes() []Axis {
	return []Axis{
		{
			Name: "Length",
			Id:   Uniq(),
			Techniques: []Technique{
				FugalistTechnique("pt.normal"),
				FugalistTechnique("pt.staccato"),
				FugalistTechnique("pt.staccatissimo"),
				FugalistTechnique("pt.tenuto"),
				FugalistTechnique("pt.portato"),
			},
			SortOrder: 0,
		},
		{
			Name: "Legato",
			Id:   Uniq(),
			Techniques: []Technique{
				FugalistTechnique("pt.normal"),
				FugalistTechnique("pt.legato"),
			},
			SortOrder: 100,
		},
		{
			Name: "Vibrato",
			Id:   Uniq(),
			Techniques: []Technique{
				FugalistTechnique("pt.normal"),
				FugalistTechnique("pt.nonVibrato"),
			},
			SortOrder: 200,
		},
		{
			Name: "Attack",
			Id:   Uniq(),
			Techniques: []Technique{
				FugalistTechnique("pt.normal"),
				FugalistTechnique("pt.marcato"),
			},
			SortOrder: 300,
		},
		{
			Name: "Technique",
			Id:   Uniq(),
			Techniques: []Technique{
				FugalistTechnique("pt.normal"),
			},
			SortOrder: 400,
		},
	}
}

func BuildOccursWith(combos []string) func(a string, b string) bool {
	mapping := make(map[string]map[string]bool)
	for _, combo := range combos {
		pts := strings.Split(combo, "+")
		for _, pt := range pts {
			m, present := mapping[pt]
			if !present {
				m = make(map[string]bool)
				mapping[pt] = m
			}
			for _, other := range pts {
				if pt != other {
					m[other] = true
				}
			}
		}
	}
	return func(a string, b string) bool {
		return mapping[a] != nil && mapping[a][b]
	}
}

// InterferesWith returns true if technique occurs with any technique in axis.
func InterferesWith(axis Axis, technique Technique, occursWith func(a string, b string) bool) bool {
	for _, tech := range axis.Techniques {
		if occursWith(tech.Id, technique.Id) {
			return true
		}
	}
	return false
}

func FindAxes(combos []string) map[AxisId]Axis {
	occursWith := BuildOccursWith(combos)

	axes := DefaultAxes()
	extras := FindExtraTechniques(combos)
	sortOrder := 100.0
outer:
	for _, extra := range extras {
		technique := FugalistTechnique(extra)
		for k := 4; k < len(axes); k++ {
			for _, t := range axes[k].Techniques {
				if t.Id == technique.Id {
					continue outer
				}
			}
			if !InterferesWith(axes[k], technique, occursWith) {
				axes[k].Techniques = append(axes[k].Techniques, technique)
				continue outer
			}
		}
		axes = append(axes, Axis{
			Name:      Uniq(),
			SortOrder: sortOrder,
			Techniques: []Technique{
				FugalistTechnique("pt.natural"),
				technique,
			},
		})
		sortOrder += 100
	}
	result := make(map[AxisId]Axis)
	for _, a := range axes {
		result[a.Id] = a
	}
	return result
}

func GetVstSounds(ptMap PtMap) map[VstSoundId]*VstSound {
	vstSounds := make(map[VstSound]bool)

	// Gather up the set of distinct (start, stop, dyn) tuples
	for _, branchMap := range ptMap {
		for _, sound := range branchMap {
			vstSound := VstSound{
				Midi:     sound.On,
				Stop:     sound.Off,
				Dynamics: sound.Dyn,
			}
			vstSounds[vstSound] = true
		}
	}

	// Create a VstSound for each distinct tuple
	vstCount := 0
	result := make(map[VstSoundId]*VstSound)
	for vstSound := range vstSounds {
		vstCount++
		vstSound.Id = Uniq()
		vstSound.Name = fmt.Sprintf("vst-%d", vstCount)
		result[vstSound.Id] = &vstSound
	}
	return result
}
