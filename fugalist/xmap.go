package fugalist

import "C"
import (
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"math"
	"sort"
	"strings"
)

func CreateExpressionMap(summary ProjectSummary, p *Project) (*doricolib.ExpressionMap, error) {
	combos, err := CreateCombinations(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create combinations: %w", err)
	}
	addOns, err := CreateTechniqueAddOns(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create addOns: %w", err)
	}

	em := doricolib.ExpressionMap{
		Name:                          summary.Name,
		EntityId:                      p.ProjectId,
		ParentEntityId:                "",
		InheritanceMask:               "0",
		Creator:                       "",
		Description:                   summary.Description,
		Version:                       fmt.Sprintf("%d", summary.Version),
		PluginNames:                   summary.Plugins,
		AutoMutualExclusion:           false,
		AllowMultipleNotesAtSamePitch: false,
		InitSwitchData: doricolib.InitSwitchData{
			Enabled: false,
			InitActions: doricolib.EntityList{
				IsArray:  "true",
				Contents: nil,
			},
		},
		Combinations:          *combos,
		TechniqueAddOns:       *addOns,
		MutualExclusionGroups: *CreateMutualExclusionGroups(p.Axes),
	}
	return &em, nil
}

func CreateCombinations(p *Project) (*doricolib.PlayingTechniqueCombinationList, error) {
	combos, err := CreateCombos(p)
	if err != nil {
		return nil, err
	}
	cl := doricolib.PlayingTechniqueCombinationList{
		IsArray: "true",
		Combos:  combos,
	}
	return &cl, nil
}

// CreateCombos creates the list of playing technique combinations.
func CreateCombos(p *Project) ([]*doricolib.PlayingTechniqueCombination, error) {
	r := make([]*doricolib.PlayingTechniqueCombination, 0)
	axes := SortedAxes(p.Axes)
	size := GetSize(axes)
	for k := 0; k < size; k++ {
		techniques, err := GetCombo(axes, k)
		if err != nil {
			return nil, err
		}

		key := GetComboKey(axes, k)

		if p.Assignments[key].Sound == "" {
			continue
		}
		soundId := p.Assignments[key].Sound

		pigment, isPigment := p.VstSounds[soundId]
		if isPigment {
			combo, err := CreateComboForPigment(techniques, pigment, p.MiddleC)
			if err != nil {
				return nil, fmt.Errorf("failed to create combo for pigment: %w", err)
			}
			r = append(r, combo)
		} else {
			color, isColor := p.CompositeSounds[soundId]
			if !isColor {
				return nil, fmt.Errorf("no sound for %s (key %s)", techniques, key)
			}
			combos, err := CreateCombosForColor(techniques, color, p, p.MiddleC)
			if err != nil {
				return nil, fmt.Errorf("failed to create combos for color: %w", err)
			}
			r = append(r, combos...)
		}
	}
	return r, nil
}

func CreateComboForPigment(techniques string, pigment *VstSound, middleC string) (*doricolib.PlayingTechniqueCombination, error) {
	volSpec, volRange, err := ParseVolumeSpec(pigment.Dynamics)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dynamics: %w", err)
	}
	switchOnActions, err := ParseSwitchOnActionList(pigment.Midi, middleC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse switch-on actions: %w", err)
	}
	combo := &doricolib.PlayingTechniqueCombination{
		TechniqueIDs:    techniques,
		BaseSwitchID:    0,
		Enabled:         true,
		Flags:           0, // TODO
		ConditionString: "",
		VelocityRange:   volRange,
		PitchRange:      "0,127",
		Transpose:       0,
		TicksBefore:     0,
		VelocityFactor:  "1.0",
		VolumeType:      *volSpec,
		AttackType: doricolib.AttackType{
			Type:   "kNoteVelocity",
			Param1: "0",
		},
		SwitchOnActions:  *switchOnActions,
		SwitchOffActions: doricolib.SwitchOffActionList{},
	}
	return combo, nil
}

func CreateCombosForColor(techniques string, color *CompositeSound, p *Project, middleC string) ([]*doricolib.PlayingTechniqueCombination, error) {
	combos := make([]*doricolib.PlayingTechniqueCombination, len(color.Branches))

	k := 0
	for _, branch := range color.Branches {
		pigment, isPigment := p.VstSounds[branch.VstSoundId]
		if !isPigment {
			return nil, fmt.Errorf("no such pigment")
		}
		combo, err := CreateComboForPigment(techniques, pigment, middleC)
		if err != nil {
			return nil, fmt.Errorf("failed to create combo for pigment: %w", err)
		}
		cond, err := Input(branch.Condition).ParseCondition()
		if err != nil {
			return nil, fmt.Errorf(`failed to parse condition: "%s"`, branch.Condition)
		}
		combo.ConditionString = cond.String()

		// Note length
		if math.IsNaN(branch.Length) {
			combo.Flags = 0
		} else {
			combo.Flags = 1
			combo.LengthFactor = fmt.Sprintf("%f", branch.Length/100.0)
		}

		// Transpose
		if math.IsNaN(branch.Transpose) {
			combo.Transpose = 0
		} else {
			combo.Transpose = int(branch.Transpose)
		}

		combos[k] = combo
		k++
	}
	return combos, nil
}

func SortedAxes(axes map[string]Axis) []Axis {
	result := make([]Axis, len(axes))
	k := 0
	for _, axis := range axes {
		result[k] = axis
		k++
	}
	sort.Slice(result, func(a, b int) bool { return result[a].SortOrder < result[b].SortOrder })
	return result
}

func GetCombo(axes []Axis, k int) (string, error) {
	result := make([]string, 0)
	for a := len(axes) - 1; a >= 0; a-- {
		axis := axes[a]
		ind := k % len(axis.Techniques)
		k = k / len(axis.Techniques)
		if ind == 0 {
			continue
		}
		technique := doricolib.GetTechniqueByName(axis.Techniques[ind].Name).Id
		result = append(result, technique)
	}
	if len(result) == 0 {
		return "pt.natural", nil
	} else {
		return strings.Join(result, "+"), nil
	}
}

func GetComboKey(axes []Axis, k int) string {
	result := make([]string, len(axes))
	for a := len(axes) - 1; a >= 0; a-- {
		axis := axes[a]
		ind := k % len(axis.Techniques)
		k = k / len(axis.Techniques)
		id := axis.Techniques[ind].Id
		result[a] = id
	}
	return Xor(result)
}

func GetSize(axes []Axis) int {
	r := 1
	for _, axis := range axes {
		r *= len(axis.Techniques)
	}
	return r
}

func CreateMutualExclusionGroups(axes map[string]Axis) *doricolib.MutexGroupList {
	groups := make([]*doricolib.MutualExclusionGroup, len(axes))
	k := 0
	for _, axis := range axes {
		if len(axis.Techniques) > 1 {
			groups[k] = MutexGroup(&axis)
			k++
		}
	}
	return &doricolib.MutexGroupList{
		IsArray:               "true",
		MutualExclusionGroups: groups,
	}
}

func MutexGroup(axis *Axis) *doricolib.MutualExclusionGroup {
	techniques := make([]string, len(axis.Techniques)-1)
	for k := 1; k < len(axis.Techniques); k++ {
		techniques[k-1] = doricolib.GetTechniqueByName(axis.Techniques[k].Name).Id
	}
	return &doricolib.MutualExclusionGroup{
		GroupId:      fmt.Sprintf("ptmg.user.%s", axis.Id),
		Name:         axis.Name,
		TechniqueIds: strings.Join(techniques, ", "),
	}
}

func rangeString(limits []string) string {
	if limits[0] == "" {
		return "0,127"
	}
	return limits[0] + "," + limits[1]
}

func CreateTechniqueAddOns(p *Project) (*doricolib.TechniqueAddOnList, error) {
	addOns := make([]doricolib.TechniqueAddOn, len(p.Tints))
	k := 0
	for _, modifier := range p.Tints {
		addOn, err := CreateTechniqueAddOn(*modifier, p.MiddleC)
		if err != nil {
			return nil, fmt.Errorf("failed to create add-on: %w", err)
		}
		addOns[k] = *addOn
		k++
	}
	result := &doricolib.TechniqueAddOnList{
		IsArray:         "true",
		TechniqueAddOns: addOns,
	}
	return result, nil
}

func CreateTechniqueAddOn(modifier Tint, middleC string) (*doricolib.TechniqueAddOn, error) {
	start, err := ParseSwitchOnActionList(modifier.Midi, middleC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start action list: %w", err)
	}
	stop, err := ParseSwitchOffActionList(modifier.Stop, middleC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop action list: %w", err)
	}
	return &doricolib.TechniqueAddOn{
		SwitchID:         0,
		TechniqueIDs:     doricolib.GetTechniqueByName(modifier.Name).Id,
		Enabled:          true,
		SwitchOnActions:  *start,
		SwitchOffActions: *stop,
	}, nil
}
