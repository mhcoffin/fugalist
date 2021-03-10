package fugalist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClause(t *testing.T) {
	tests := []struct {
		name     string
		err      bool
		expected Clause
	}{
		{"NoteLength < long", false, Clause{operator: LT, lhs: NoteLength, rhs: Long}},
		{"nl == short", false, Clause{operator: EQ, lhs: NoteLength, rhs: Short}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, clause, err := ParseClause(Input(test.name))
			if test.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.expected, clause)
			}

		})
	}
}

func TestInput_ParseCondition(t *testing.T) {
	tests := []struct {
		name     string
		expected Condition
		ok       bool
	}{
		{"nl < short", Condition{connector: NoConjunction, clauses: []Clause{{operator: LT, lhs: NoteLength, rhs: Short}}}, true},
		{"NoteLength < short", Condition{connector: NoConjunction, clauses: []Clause{{operator: LT, lhs: NoteLength, rhs: Short}}}, true},
		{"nl < long and nl >= VeryShort",
			Condition{
				connector: And,
				clauses: []Clause{
					{operator: LT, lhs: NoteLength, rhs: Long},
					{operator: GE, lhs: NoteLength, rhs: VeryShort},
				}},
			true,
		},
		{"nl < long and nl >= VeryShort and nl != m",
			Condition{
				connector: And,
				clauses: []Clause{
					{operator: LT, lhs: NoteLength, rhs: Long},
					{operator: GE, lhs: NoteLength, rhs: VeryShort},
					{operator: NE, lhs: NoteLength, rhs: Medium},
				}},
			true,
		},
		{"nl < long || nl >= VeryShort || nl != m",
			Condition{
				connector: Or,
				clauses: []Clause{
					{operator: LT, lhs: NoteLength, rhs: Long},
					{operator: GE, lhs: NoteLength, rhs: VeryShort},
					{operator: NE, lhs: NoteLength, rhs: Medium},
				}},
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond, err := Input(test.name).ParseCondition()
			assert.Equalf(t, test.ok, err == nil, "%s", err)
			assert.Equal(t, test.expected, cond)
		})
	}
}

func TestCondition_String(t *testing.T) {
	tests := []struct {
		name string
		expected string
	}{
		{"NoteLength < short", "NoteLength &LT; kShort"},
		{"nl < s", "NoteLength &LT; kShort"},
		{"nl > s & nl <= vl", "NoteLength &GT; kShort AND NoteLength &LT;= kVeryLong"},
		{"nl == s | nl == vs", "NoteLength == kShort OR NoteLength == kVeryShort"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond, err := Input(test.name).ParseCondition()
			assert.Nil(t, err)
			assert.Equal(t, test.expected, cond.String())
		})
	}
}

