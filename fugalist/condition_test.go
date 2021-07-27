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
		{"short < note length", false, Clause{operator: GT, lhs: NoteLength, rhs: Short}},
		{"note length < long", false, Clause{operator: LT, lhs: NoteLength, rhs: Long}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, clause, err := Input(test.name).ParseClause()
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
		{"nl < short", Condition{connector: And, clauses: []Clause{{operator: LT, lhs: NoteLength, rhs: Short}}}, true},
		{"NoteLength < short", Condition{connector: And, clauses: []Clause{{operator: LT, lhs: NoteLength, rhs: Short}}}, true},
		{"nl < long and nl >= VeryShort",
			Condition{
				connector: And,
				clauses: []Clause{
					{operator: LT, lhs: NoteLength, rhs: Long},
					{operator: GE, lhs: NoteLength, rhs: VeryShort},
				}},
			true,
		},
		{"nl < long AND nl >= VeryShort AND nl != medium",
			Condition{
				connector: And,
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
			cond, err := Input(test.name).ParseClauseList()
			assert.Equalf(t, test.ok, err == nil, "%s", err)
			assert.Equal(t, test.expected, cond)
		})
	}
}

func TestClauseList_String(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"NoteLength < short", "NoteLength < kShort"},
		{"nl < short", "NoteLength < kShort"},
		{"nl > short AND nl <= veryLong", "NoteLength > kShort AND NoteLength <= kVeryLong"},
		{"nl == Short and nl == very short", "NoteLength == kShort AND NoteLength == kVeryShort"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond, err := Input(test.name).ParseClauseList()
			assert.Nil(t, err)
			assert.Equal(t, test.expected, cond.String())
		})
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		expected Condition
	}{
		{
			"short < NoteLength < long", Condition{
				connector: And,
				clauses: []Clause{
					{
						operator: GT,
						lhs:      "NoteLength",
						rhs:      "short",
					},
					{
						operator: LT,
						lhs:      "NoteLength",
						rhs:      "long",
					},
				},
			},
		},
		{
			"veryShort <= NoteLength <= veryLong", Condition{
				connector: And,
				clauses: []Clause{
					{
						operator: GE,
						lhs:      "NoteLength",
						rhs:      "veryShort",
					},
					{
						operator: LE,
						lhs:      "NoteLength",
						rhs:      "veryLong",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cond, err := Input(test.name).ParseRange()
			assert.Nil(t, err)
			assert.Equal(t, test.expected, cond)
		})
	}
}

func TestRange_Failure(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"veryShort >= veryLong"},
		{"veryShort <= NoteLength > long"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Input(test.name).ParseRange()
			assert.NotNil(t, err)
		})
	}
}

func TestInput_ParseBranch(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{" very short < NoteLength ", "NoteLength > kVeryShort"},
		{"short < note length < long", "NoteLength > kShort AND NoteLength < kLong"},
		{"short < note length and note length < long", "NoteLength > kShort AND NoteLength < kLong"},
		{" very   short < note length and note length < veryLong ", "NoteLength > kVeryShort AND NoteLength < kVeryLong"},
		{"short <= noteLength < medium", "NoteLength >= kShort AND NoteLength < kMedium"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run(test.name, func(t *testing.T) {
				b, err := Input(test.name).ParseBranch()
				assert.Nil(t, err)
				assert.Equal(t, test.expected, b.String())
			})
		})
	}
}
