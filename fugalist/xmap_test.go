package fugalist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var lenAxis = Axis{
	Id:   "abc",
	Name: "Length",
	Techniques: []Technique{
		{"abc2", "Normal"},
		{"abc1", "Staccato"},
		{"abc3", "Tenuto"},
	},
}

var techAxis = Axis{
	Id:   "Tech",
	Name: "Techniques",
	Techniques: []Technique{
		{"Tech1", "Normal"},
		{"Tech2", "Pizzicato"},
		{"Tech3", "Flautando"},
	},
}

func Test(t *testing.T) {
	tests := []struct {
		name     string
		ind      int
		expected string
	}{
		{"0", 0, "pt.natural"},
		{"1", 1, "pt.pizzicato"},
		{"2", 2, "pt.flautando"},
		{"3", 3, "pt.staccato"},
		{"4", 4, "pt.pizzicato+pt.staccato"},
		{"5", 5, "pt.flautando+pt.staccato"},
		{"6", 6, "pt.tenuto"},
		{"7", 7, "pt.pizzicato+pt.tenuto"},
		{"8", 8, "pt.flautando+pt.tenuto"},
	}
	axes := []Axis{lenAxis, techAxis}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			combo, err := GetCombinationString(axes, test.ind)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, combo)
		})
	}
}
