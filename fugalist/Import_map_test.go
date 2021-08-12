package fugalist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanonicalizeTechniqueString(t *testing.T) {
	tests := []struct {
		name string
		orig string
		expected string
	}{
		{"single", "pt.legato", "pt.legato"},
		{"sorted", "pt.legato+pt.nonVibrato", "pt.legato+pt.nonVibrato"},
		{"reversed", "pt.nonVibrato+pt.legato", "pt.legato+pt.nonVibrato"},

	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, CanonicalizeTechniqueString(test.orig))
		})
	}
}
