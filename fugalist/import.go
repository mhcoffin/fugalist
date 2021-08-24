package fugalist

import (
	"fmt"
	"github.com/mhcoffin/go-doricolib/doricolib"
	"strings"
)

func ImportDynamics(volType doricolib.VolumeType, volRange string) (string, error) {
	rng := ""
	if volRange != "" {
		r := strings.Split(volRange, ",")
		if len(r) != 2 {
			return "", fmt.Errorf("bad volume Range: %s", volRange)
		}
		rng = fmt.Sprintf(" %s:%s", r[0], r[1])
	}
	switch volType.Type {
	case "kCC":
		return fmt.Sprintf("CC%s%s", volType.Param1, rng), nil
	default:
		return fmt.Sprintf("velocity%s", rng), nil
	}
}

func ImportSwitchOnActions(switchActions doricolib.SwitchOnActionList) (string, error) {
	result := make([]string, len(switchActions.SwitchOnActions))
	for k, action := range switchActions.SwitchOnActions {
		switch action.Type {
		case "kKeySwitch":
			vel := ""
			if action.Param2 != "127" {
				vel = fmt.Sprintf("=%s", action.Param2)
			}
			result[k] = fmt.Sprintf("KS%s%s", action.Param1, vel)
		case "kControlChange":
			result[k] = fmt.Sprintf("CC%s=%s", action.Param1, action.Param2)
		case "kProgramChange":
			result[k] = fmt.Sprintf("PC%s", action.Param1)
		}
	}
	return strings.Join(result, ", "), nil
}
