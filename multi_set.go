package gotestrail

type MultiSet struct {
	CaseSet     *CaseSet     `json:"cases"`
	CaseTypeSet *CaseTypeSet `json:"caseTypes"`
	SectionSet  *SectionSet  `json:"sections"`
}

func NewMultiSet() *MultiSet {
	return &MultiSet{
		CaseSet:     NewCaseSet(),
		CaseTypeSet: &CaseTypeSet{},
		SectionSet:  NewSectionSet(),
	}
}

func (set *MultiSet) Lens() map[string]uint {
	m := map[string]uint{}
	if set.CaseSet != nil {
		m[SlugCase] = set.CaseSet.Len()
	}
	if set.SectionSet != nil {
		m[SlugSection] = set.SectionSet.Len()
	}
	return m
}
