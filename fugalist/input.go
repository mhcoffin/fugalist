package fugalist

import (
	"fmt"
	"regexp"
	"unicode"
)

type ComparisonOperator int

const (
	NoComparison ComparisonOperator = iota
	LT
	LE
	EQ
	NE
	GT
	GE
)

func (com ComparisonOperator) String() string {
	switch com {
	case LT:
		return "<"
	case LE:
		return "<="
	case EQ:
		return "=="
	case NE:
		return "!="
	case GT:
		return ">"
	case GE:
		return ">="
	default:
		panic("no such comparison operator")
	}
}

func (com ComparisonOperator) Opposite() ComparisonOperator {
	switch com {
	case LT:
		return GT
	case LE:
		return GE
	case EQ:
		return EQ
	case NE:
		return NE
	case GT:
		return LT
	case GE:
		return GT
	default:
		panic("no such comparison operator")
	}
}

var ComparisonOperatorMap = map[string]ComparisonOperator{
	"":   NoComparison,
	"<=": LE,
	"<":  LT,
	">=": GE,
	">":  GT,
	"==": EQ,
	"!=": NE,
}

type Constant string

const (
	NoConstant Constant = "?"
	VeryShort  Constant = "veryShort"
	Short      Constant = "short"
	Medium     Constant = "medium"
	Long       Constant = "long"
	VeryLong   Constant = "veryLong"
)

func (con Constant) String() string {
	switch con {
	case VeryShort:
		return "kVeryShort"
	case Short:
		return "kShort"
	case Medium:
		return "kMedium"
	case Long:
		return "kLong"
	case VeryLong:
		return "kVeryLong"
	default:
		panic("no such constant")
	}
}

type Variable string

const (
	NoVariable Variable = "?"
	NoteLength Variable = "NoteLength"
)

func (v Variable) String() string {
	switch v {
	case NoteLength:
		return "NoteLength"
	default:
		panic("not a variable")
	}
}

type Conjunction int

const (
	NoConjunction Conjunction = iota
	And
)

func (conj Conjunction) String() string {
	return "AND"
}

type Input string

func (in Input) String() string {
	return string(in)
}

func (in Input) Empty() bool {
	s := in.SkipWhitespace()
	return len(s) == 0
}

func (in Input) SkipWhitespace() Input {
	for k, c := range in {
		if !unicode.IsSpace(c) {
			return in[k:]
		}
	}
	return ""
}

func (in Input) MustBe(regexp *regexp.Regexp) (Input, string, error) {
	inp := in.SkipWhitespace()
	pos := regexp.FindStringIndex(string(inp))
	if pos == nil {
		return inp, "", fmt.Errorf("expected to match %s at %s", regexp.String(), inp.String())
	}
	return inp[pos[1]:], string(inp[pos[0]:pos[1]]), nil
}

func (in Input) MustBeIdentifier() (Input, string, error) {
	r := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*`)
	r.Longest()
	rest, id, err := in.MustBe(r)
	if err != nil {
		return rest, id, fmt.Errorf(`identifier expected: "%s"`, in)
	}
	return rest, id, err
}

var noteLengthRegexp = regexp.MustCompile(`(?i)(note\s*length)|(nl)`)

func (in Input) MustBeVariable() (Input, Variable, error) {
	rest, _, err := in.MustBe(noteLengthRegexp)
	if err != nil {
		return in, NoVariable, fmt.Errorf("variable name expected: %w", err)
	}
	return rest, "NoteLength", nil
}

var veryShortRegexp = regexp.MustCompile(`^(?i)(very\s*short)`)
var shortRegexp = regexp.MustCompile(`^(?i)short`)
var mediumRegexp = regexp.MustCompile(`^(?i)medium`)
var longRegexp = regexp.MustCompile(`^(?i)long`)
var veryLongRegexp = regexp.MustCompile(`^(?i)very\s*long`)

func (in Input) MustBeConstant() (Input, Constant, error) {
	rest, _, err := in.MustBe(veryShortRegexp)
	if err == nil {
		return rest, "veryShort", nil
	}
	rest, _, err = in.MustBe(shortRegexp)
	if err == nil {
		return rest, "short", nil
	}
	rest, _, err = in.MustBe(mediumRegexp)
	if err == nil {
		return rest, "medium", nil
	}
	rest, _, err = in.MustBe(longRegexp)
	if err == nil {
		return rest, "long", nil
	}
	rest, _, err = in.MustBe(veryLongRegexp)
	if err == nil {
		return rest, "veryLong", nil
	}
	return rest, "?", fmt.Errorf("expected constant")
}

func (in Input) MustBeComparisonOperator() (Input, ComparisonOperator, error) {
	r := regexp.MustCompile(`^(<|<=|>|>=|==|!=)`)
	r.Longest()
	rest, op, err := in.MustBe(r)
	if err != nil {
		return in, NoComparison, err
	}
	return rest, ComparisonOperatorMap[op], nil
}

func (in Input) MustBeConjunction() (Input, Conjunction, error) {
	r := regexp.MustCompile(`(?i)^(and)`)
	r.Longest()

	rest, _, err := in.MustBe(r)
	if err != nil {
		return in, NoConjunction, err
	}
	return rest, And, nil
}
