package fugalist

import (
	"fmt"
	"regexp"
	"strings"
)

type Clause struct {
	operator ComparisonOperator
	lhs      Variable
	rhs      Constant
}

func (clause *Clause) String() string {
	return fmt.Sprintf("%s %s %s", clause.lhs, clause.operator, clause.rhs)
}

type Condition struct {
	connector Conjunction
	clauses   []Clause
}

func (cond *Condition) String() string {
	clauses := make([]string, len(cond.clauses))
	for k, clause := range cond.clauses {
		clauses[k] = clause.String()
	}
	conj := fmt.Sprintf(" %s ", cond.connector.String())
	return strings.Join(clauses, conj)
}

func (in Input) ParseBranch() (Condition, error) {
	cond, err := in.ParseRange()
	if err == nil {
		return cond, nil
	}
	cond, err = in.ParseClauseList()
	return cond, err
}

func (in Input) ParseClauseList() (Condition, error) {
	connector := NoConjunction
	clauses := make([]Clause, 0)
	var rest = in
	var err error
	for !rest.Empty() {
		var clause Clause
		rest, clause, err = rest.ParseClause()
		if err != nil {
			return Condition{}, err
		}
		clauses = append(clauses, clause)
		if !rest.Empty() {
			var c Conjunction
			rest, c, err = rest.MustBeConjunction()
			if err != nil {
				return Condition{}, fmt.Errorf("AND or OR expected")
			}
			if connector != NoConjunction && c != connector {
				return Condition{}, fmt.Errorf("inconsistent AND/OR combination")
			}
			connector = c
		}
	}
	return Condition{connector, clauses}, nil
}

func (in Input) ParseComparisonNoteLengthFirst() (Input, Clause, error) {
	var rest = in
	var result = Clause{}
	var err error
	rest, result.lhs, err = rest.MustBeVariable()
	if err != nil {
		return in, Clause{}, fmt.Errorf("identifier expected")
	}
	rest, result.operator, err = rest.MustBeComparisonOperator()
	if err != nil {
		return in, Clause{}, fmt.Errorf("comparison operator expected")
	}
	rest, result.rhs, err = rest.MustBeConstant()
	if err != nil {
		return in, Clause{}, fmt.Errorf("constant expected")
	}
	return rest, result, nil
}

func (in Input) ParseComparisonLengthConstantFirst() (Input, Clause, error) {
	var rest = in
	var result = Clause{}
	var err error
	rest, result.rhs, err = rest.MustBeConstant()
	if err != nil {
		return in, Clause{}, fmt.Errorf("identifier expected")
	}
	rest, result.operator, err = rest.MustBeComparisonOperator()
	if err != nil {
		return in, Clause{}, fmt.Errorf("comparison operator expected")
	}
	result.operator = result.operator.Opposite()
	rest, result.lhs, err = rest.MustBeVariable()
	if err != nil {
		return in, Clause{}, fmt.Errorf("constant expected")
	}
	return rest, result, nil
}

func (in Input) ParseClause() (Input, Clause, error) {
	rest, result, err := in.ParseComparisonNoteLengthFirst()
	if err == nil {
		return rest, result, err
	}
	rest, result, err = in.ParseComparisonLengthConstantFirst()
	if err == nil {
		return rest, result, err
	}
	return in, Clause{}, fmt.Errorf("clause expected")
}

func (in Input) ParseRange() (Condition, error) {
	var rest = in
	var lhs Constant
	var op1 string
	var nl Variable
	var op2 string
	var rhs Constant
	var err error

	rest, lhs, err = rest.MustBeConstant()
	if err != nil {
		return Condition{}, fmt.Errorf("length constant expected")
	}
	rest, op1, err = rest.MustBe(regexp.MustCompile("<=?"))
	if err != nil {
		return Condition{}, fmt.Errorf("range operator expected")
	}
	rest, nl, err = rest.MustBeVariable()
	if err != nil {
		return Condition{}, fmt.Errorf("NoteLength expected")
	}
	rest, op2, err = rest.MustBe(regexp.MustCompile("<=?"))
	if err != nil {
		return Condition{}, fmt.Errorf("range operator expected")
	}
	rest, rhs, err = rest.MustBeConstant()
	if err != nil {
		return Condition{}, fmt.Errorf("length constant expected")
	}
	if op1 == "<=" {
		op1 = ">="
	} else {
		op1 = ">"
	}
	condition := Condition{
		connector: And,
		clauses: []Clause{
			{
				operator: ComparisonOperatorMap[op1],
				lhs:      nl,
				rhs:      lhs,
			},
			{
				operator: ComparisonOperatorMap[op2],
				lhs:      nl,
				rhs:      rhs,
			},
		},
	}
	return condition, nil
}
