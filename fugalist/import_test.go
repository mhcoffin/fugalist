package fugalist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImportVelocity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"ParseTest1"},
		{"Ref"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scoreLib := ReadDoricolib(t, test.name)
			for _, combo := range scoreLib.ExpressionMaps.Entities.Contents[0].Combinations.Combos {
				dyn, err := ImportDynamics(combo.VolumeType, combo.VelocityRange)
				assert.Nil(t, err)
				volType, rng, err := ParseVolumeSpec(dyn)
				assert.Nil(t, err)
				assert.Equal(t, &combo.VolumeType, volType)
				assert.Equal(t, combo.VelocityRange, rng)
			}
		})
	}
}

func TestImportSwitchOnActions(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"ParseTest1"},
		{"Ref"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scoreLib := ReadDoricolib(t, test.name)
			for _, combo := range scoreLib.ExpressionMaps.Entities.Contents[0].Combinations.Combos {
				switchOnActionList := combo.SwitchOnActions
				midi, err := ImportSwitchOnActions(switchOnActionList)
				assert.Nil(t, err)
				switchOnActions, err := ParseSwitchOnActionList(midi, "4")
				assert.Nil(t, err)
				assert.Equal(t, &switchOnActionList, switchOnActions)
			}
		})
	}
}
