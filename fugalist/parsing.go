package fugalist

import (
	"errors"
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func ParseSwitchOnActionList(s string, middleC string) (*doricolib.SwitchOnActionList, error) {
	actions, err := ParseActionList(s, middleC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start-action list: %w", err)
	}
	return &doricolib.SwitchOnActionList{
		IsArray:         "true",
		SwitchOnActions: actions,
	}, nil
}

func ParseSwitchOffActionList(s string, middleC string) (*doricolib.SwitchOffActionList, error) {
	actions, err := ParseActionList(s, middleC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop-action list: %w", err)
	}
	return &doricolib.SwitchOffActionList{
		IsArray:          "true",
		SwitchOffActions: actions,
	}, nil
}

func ParseActionList(s string, middleC string) ([]doricolib.SwitchAction, error) {
	c := 4
	switch middleC {
	case "C3", "c3":
		c = 3
	case "C4", "c4":
		c = 4
	case "C5", "c5":
		c = 5
	}
	parts := strings.Split(s, ",")
	actions := make([]doricolib.SwitchAction, 0)
	for _, part := range parts {
		action, err := ParseMidi(part, c)
		if err != nil {
			return nil, err
		}
		if action != nil {
			actions = append(actions, *action)
		}
	}
	return actions, nil
}

/*
A midi action can be one of
  CC n = m
  KS n [= m]
  PC n
  [A-G]#? n
*/
var EmptyPat = regexp.MustCompile(`^\s*$`)
var CcPat = regexp.MustCompile(`^\s*(?i:CC)\s*(\d+)\s*=\s*(\d+)\s*$`)
var CCPartPat = regexp.MustCompile(`^\s*(?i:CC)\s*(\d+)\s*=\s*(\d+)\/(\d+)\s*$`)
var KsPat = regexp.MustCompile(`^\s*(?i:KS)\s*(\d+)\s*(?:=\s*(\d+))?\s*$`)
var PcPat = regexp.MustCompile(`^\s*(?i:PC)\s*(\d+)\s*$`)
var NotePat = regexp.MustCompile(`^\s*([A-Ga-g])([#b]?)\s*(-?\d+)\s*$`)

func ParseMidi(part string, middleCOctave int) (*doricolib.SwitchAction, error) {
	switch {
	case EmptyPattern.MatchString(part):
		return nil, nil
	case CcPat.MatchString(part):
		x := CcPat.FindStringSubmatch(part)
		return &doricolib.SwitchAction{
			Type:   "kControlChange",
			Param1: x[1],
			Param2: x[2],
		}, nil
	case CCPartPat.MatchString(part):
		x := CCPartPat.FindStringSubmatch(part)
		setting, err := proportion(x[2], x[3])
		if err != nil {
			return nil, fmt.Errorf("bad CC")
		}
		return &doricolib.SwitchAction{
			Type:   "kControlChange",
			Param1: x[1],
			Param2: setting,
		}, nil
	case KsPat.MatchString(part):
		x := KsPat.FindStringSubmatch(part)
		vel := "127"
		if x[2] != "" {
			vel = x[2]
		}
		return &doricolib.SwitchAction{
			Type:   "kKeySwitch",
			Param1: x[1],
			Param2: vel,
		}, nil
	case PcPat.MatchString(part):
		x := PcPat.FindStringSubmatch(part)
		return &doricolib.SwitchAction{
			Type:   "kProgramChange",
			Param1: x[1],
			Param2: "0",
		}, nil
	case NotePat.MatchString(part):
		x := NotePat.FindStringSubmatch(part)
		number, err := note(x[1], x[2], x[3], middleCOctave)
		if err != nil {
			return nil, fmt.Errorf("failed to parse midi: %v", part)
		}
		return &doricolib.SwitchAction{
			Type:   "kKeySwitch",
			Param1: strconv.Itoa(number),
			Param2: "127",
		}, nil

	default:
		return nil, fmt.Errorf("illegal midi setting: \"%v\"", part)
	}
}

func proportion(num string, den string) (string, error) {
	n, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return "", fmt.Errorf("not a number: %s", num)
	}
	d, err := strconv.ParseFloat(den, 64)
	if err != nil {
		return "", fmt.Errorf("not a number: %s", den)
	}
	if n < 1 {
		return "", fmt.Errorf("midi numerator must be >= 1: %d/%d", int(n), int(d))
	}
	if n > d {
		return "", fmt.Errorf("fraction is not proper: %d/%d", int(n), int(d))
	}
	value := math.Round(((n - 0.5) / d) * 128.0)
	return fmt.Sprintf("%d", int(value)), nil

}

func note(name, sharpOrFlat, octave string, middleCOctave int) (int, error) {
	noteNumber, err := noteNumber(name)
	if err != nil {
		return -1, err
	}
	switch sharpOrFlat {
	case "#":
		noteNumber++
	case "b":
		noteNumber--
	}

	oct, err := strconv.Atoi(octave)
	if err != nil {
		return -1, err
	}
	return 60 + (oct-middleCOctave)*12 + noteNumber, nil
}

func noteNumber(noteName string) (int, error) {
	switch noteName {
	case "C", "c":
		return 0, nil
	case "D", "d":
		return 2, nil
	case "E", "e":
		return 4, nil
	case "F", "f":
		return 5, nil
	case "G", "g":
		return 7, nil
	case "A", "a":
		return 9, nil
	case "B", "b":
		return 11, nil
	default:
		return -1, fmt.Errorf("bad note name: %s", noteName)
	}
}

var transposePattern = regexp.MustCompile(`^\s+([+-]?)(\d+)\s*$`)

func ParseTranspose(t string) (int, error) {
	if !transposePattern.MatchString(t) {
		return 0, errors.New("illegal transpose")
	}
	parts := transposePattern.FindStringSubmatch(t)
	v, err := strconv.Atoi(parts[0] + parts[1])
	if err != nil {
		return 0, err
	}
	return v, nil
}

//
// The volume spec looks like
//   CCn min:mix
// or
//   Velocity min:max
//
var EmptyPattern = regexp.MustCompile(`^\s*$`)
var CcPattern = regexp.MustCompile(`^\s*(?i:cc)\s*(\d+)(?:\s+(?:(\d+)\s*:\s*(\d+)))?\s*$`)
var VelPattern = regexp.MustCompile(`^\s*(?i:velocity)\s*(?:\s+(?:(\d+)\s*:\s*(\d+)))?\s*$`)

func ParseVolumeSpec(s string) (*doricolib.VolumeType, string, error) {
	switch {
	case EmptyPattern.MatchString(s):
		return &doricolib.VolumeType{
			Type:   "kNoteVelocity",
			Param1: "0",
		}, "0,127", nil
	case CcPattern.MatchString(s):
		parts := CcPattern.FindStringSubmatch(s)
		return &doricolib.VolumeType{
			Type:   "kCC",
			Param1: parts[1],
		}, rangeString(parts[2:]), nil
	case VelPattern.MatchString(s):
		parts := VelPattern.FindStringSubmatch(s)
		return &doricolib.VolumeType{
			Type:   "kNoteVelocity",
			Param1: "0",
		}, rangeString(parts[1:]), nil
	default:
		return nil, "", errors.New("bad velocity pattern")
	}
}
