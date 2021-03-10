package fugalistdb

type AssignmentRow = struct {
	Techniques []TechniqueId
	Result     Assignment
	Duplicate  bool
}

type Table []AssignmentRow

func CreateAssignmentTable(p *Project) (Table, error) {
	axes := SortedAxes(p.Axes)
	size := GetSize(axes)
	rows := make([]AssignmentRow, size)
	for k := 0; k < size; k++ {
		techniqueIds := GetTechniqueIds(axes, k)
		key := Xor(techniqueIds)
		sound := p.Assignments[key]
		rows[k] = AssignmentRow{
			Techniques: GetTechniqueIds(axes, k),
			Result:     sound,
			Duplicate:  false,
		}
	}
	return rows, nil
}

func GetTechniqueIds(axes []Axis, k int) []TechniqueId {
	result := make([]TechniqueId, len(axes))
	for a := len(axes) - 1; a >= 0; a-- {
		axis := axes[a]
		if axis.AddOn {
			continue
		}
		ind := k % len(axis.Techniques)
		k = k / len(axis.Techniques)
		result[a] = axis.Techniques[ind].Id
	}
	return result
}
