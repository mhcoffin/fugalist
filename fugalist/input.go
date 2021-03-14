package fugalist

import (
	"fmt"
	"regexp"
	"strings"
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

var ConstantMap = map[string]Constant{
	"veryshort": VeryShort,
	"vs":        VeryShort,
	"short":     Short,
	"s":         Short,
	"medium":    Medium,
	"m":         Medium,
	"long":      Long,
	"l":         Long,
	"verylong":  VeryLong,
	"vl":        VeryLong,
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

var VariableMap = map[string]Variable{
	"notelength": NoteLength,
	"nl":         NoteLength,
}

type Conjunction int

const (
	NoConjunction Conjunction = iota
	And
	Or
)

func (c Conjunction) String() string {
	switch c {
	case And:
		return "AND"
	case Or:
		return "OR"
	case NoConjunction:
		return ""
	default:
		panic("no such conjunction")
	}
}

var ConjunctionMap = map[string]Conjunction{
	"and": And,
	"&":   And,
	"&&":  And,
	"or":  Or,
	"|":   Or,
	"||":  Or,
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

func (in Input) MustBeVariable() (Input, Variable, error) {
	rest, id, err := in.MustBeIdentifier()
	if err != nil {
		return in, NoVariable, fmt.Errorf("variable name expected: %w", err)
	}
	v := VariableMap[strings.ToLower(id)]
	if v == "" {
		return in, NoVariable, fmt.Errorf("unrecognized variable: %s", id)
	}
	return rest, v, nil
}

func (in Input) MustBeConstant() (Input, Constant, error) {
	rest, id, err := in.MustBeIdentifier()
	if err != nil {
		return in, NoConstant, fmt.Errorf("constant expected: %w", err)
	}
	c := ConstantMap[strings.ToLower(id)]
	if c == "" {
		return in, NoConstant, fmt.Errorf("unknown constant: %s", id)
	}
	return rest, c, nil
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
	r := regexp.MustCompile(`^(?i:and\b)|(?i:or\b)|&|&&|(\|)|(\|\|)`)
	r.Longest()

	rest, op, err := in.MustBe(r)
	if err != nil {
		return in, NoConjunction, err
	}
	return rest, ConjunctionMap[strings.ToLower(op)], nil
}