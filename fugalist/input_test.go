package fugalist

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInput_String(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected string
	}{
		{"empty", "", ""},
		{"one", "a", "a"},
		{"longer", "superman", "superman"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inp := Input(test.s)
			assert.Equal(t, test.expected, inp.String())
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name     string
		x        Input
		expected bool
	}{
		{"empty", Input(""), true},
		{"not empty", Input("abcd"), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.x.Empty(), test.expected)
		})
	}
}

func TestInput_SkipWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		in       Input
		expected Input
	}{
		{"empty", Input(""), Input("")},
		{"no white space", Input("xyz"), Input("xyz")},
		{"no white space at beginning", Input("xyz  "), Input("xyz  ")},
		{"white space at beginning", Input("  xyz"), Input("xyz")},
		{"white space middle", Input("xyz  def"), Input("xyz  def")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.in.SkipWhitespace())
		})
	}
}

func TestInput_MustBe(t *testing.T) {
	tests := []struct {
		name     string
		in       Input
		regex    *regexp.Regexp
		expected string
		rest     Input
	}{
		{"match", Input("abcde fghi"), regexp.MustCompile(`^abcde`), "abcde", Input(" fghi")},
		{"match", Input("abcde fghi"), regexp.MustCompile(`^a?`), "a", Input("bcde fghi")},
		{"match", Input("abcde fghi"), regexp.MustCompile(`^[abcd]+`), "abcd", Input("e fghi")},
		{"match", Input("abccdde fghi"), regexp.MustCompile(`^[abcd]+`), "abccdd", Input("e fghi")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rest, actual, err := test.in.MustBe(test.regex)
			assert.Nil(t, err)
			assert.Equal(t, test.rest, rest)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestInput_MustBe_Fails(t *testing.T) {
	tests := []struct {
		name  string
		in    Input
		regex *regexp.Regexp
	}{
		{"empty", "", regexp.MustCompile("^xyzzy")},
		{"not empty", "foobar", regexp.MustCompile("^bar")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x, m, err := test.in.MustBe(test.regex)
			assert.NotNil(t, err)
			assert.Equal(t, "", m)
			assert.Equal(t, x, test.in)
		})
	}
}

func TestInput_MustBeIdentifier(t *testing.T) {
	tests := []struct {
		name string
		in   Input
		id   string
		rest Input
	}{
		{"ok", "abc def", "abc", Input(" def")},
		{"ok", " abc123<4", "abc123", Input("<4")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rest, id, err := test.in.MustBeIdentifier()
			assert.Nil(t, err)
			assert.Equal(t, test.id, id)
			assert.Equal(t, test.rest, rest)
		})
	}
}

func TestInput_MustBeComparisonOperator(t *testing.T) {
	tests := []struct {
		name     string
		in       Input
		expected ComparisonOperator
	}{
		{name: "<", in: Input("  < "), expected: LT},
		{name: "<=", in: Input("  <="), expected: LE},
		{name: ">", in: Input("  > "), expected: GT},
		{name: ">=", in: Input("  >= "), expected: GE},
		{name: "==", in: Input("  == "), expected: EQ},
		{name: "!=", in: Input("!=foobar"), expected: NE},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, op, err := test.in.MustBeComparisonOperator()
			assert.Nil(t, err)
			assert.Equal(t, test.expected, op)
		})
	}
}

func TestInput_MustBeVariable(t *testing.T) {
	tests := []struct {
		name string
		in   Input
		v    Variable
		err  bool
	}{
		{"empty", "  ", NoVariable, true},
		{"long form", " NoteLength ", NoteLength, false},
		{"short form", " nl ", NoteLength, false},
		{"case", " noteLength ", NoteLength, false},
		{"no", " note ", NoVariable, true},
		{"no", " ", NoVariable, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, variable, err := test.in.MustBeVariable()
			assert.Equal(t, test.v, variable)
			if test.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestInput_MustBeConstant(t *testing.T) {
	tests := []struct {
		name string
		in   Input
		c    Constant
		err  bool
	}{
		{"empty", "  ", NoConstant, true},
		{"very short", " VeryShort ", VeryShort, false},
		{"very short 2", " vs ", VeryShort, false},
		{"short", " short ", Short, false},
		{"short 2", " s ", Short, false},
		{"medium", " medium ", Medium, false},
		{"medium 2", " m ", Medium, false},
		{"long", "Long", Long, false},
		{"long 2", "  l foo", Long, false},
		{"no", "xyz", NoConstant, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, constant, err := test.in.MustBeConstant()
			assert.Equal(t, test.c, constant)
			if test.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestInput_MustBeConjunction(t *testing.T) {
	tests := []struct {
		name     string
		expected Conjunction
		ok       bool
	}{
		{"and", And, true},
		{" &", And, true},
		{" && NL < short", And, true},
		{" AND NL < short", And, true},
		{" ANDOVER NL < short", NoConjunction, false},
		{" or NL < short", Or, true},
		{" Or NL < short", Or, true},
		{" | NL < short", Or, true},
		{" || NL < short", Or, true},
		{" orlon NL < short", NoConjunction, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, conj, err := Input(test.name).MustBeConjunction()
			assert.Equal(t, test.expected, conj)
			if test.ok {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
