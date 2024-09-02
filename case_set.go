package gotestrail

import (
	"encoding/json"
	"os"
	"slices"
	"strings"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
)

type Case struct {
	ID                   uint    `json:"id"`
	Title                string  `json:"title"`
	CreatedBy            uint    `json:"created_by"`
	CreatedOn            uint    `json:"created_on"`
	CustomAutomationType *int    `json:"custom_automation_type"`
	DisplayOrder         *int    `json:"display_order"`
	Estimate             *string `json:"estimate"`
	EstimateForecast     *string `json:"estimate_forecast"`
	UpdatedBy            uint    `json:"updated_by"`
	UpdatedOn            uint    `json:"updated_on"`
	IsDeleted            *int    `json:"is_deleted"`
	MilestoneID          *uint   `json:"milestone_id"`
	PriorityID           *uint   `json:"priority_id"`
	Refs                 *string `json:"refs"`
	SectionID            *uint   `json:"section_id"`
	SuiteID              *uint   `json:"suite_id"`
	TemplateID           *uint   `json:"template_id"`
	TypeID               *uint   `json:"type_id"`
}

type FuncCaseMatch func(c Case) bool

func (c Case) MatchFunc(fn FuncCaseMatch) bool {
	if fn == nil {
		return false
	} else {
		return fn(c)
	}
}

func (c Case) RefsContains(s string) bool {
	if c.Refs == nil {
		return false
	} else {
		return strings.Contains(pointer.Dereference(c.Refs), s)
	}
}

type CaseSet struct {
	Cases map[uint]Case `json:"cases"`
}

func NewCaseSet() *CaseSet {
	return &CaseSet{Cases: map[uint]Case{}}
}

func ReadFileCaseSet(filename string) (*CaseSet, error) {
	set := NewCaseSet()
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return set, json.Unmarshal(b, set)
	}
}

func (set *CaseSet) ReadFileJSON(filename string) error {
	if new, err := ReadFileCaseSet(filename); err != nil {
		return err
	} else {
		for k, v := range new.Cases {
			set.Cases[k] = v
		}
	}
	return nil
}

func (set *CaseSet) Add(c ...Case) {
	for _, ci := range c {
		set.Cases[ci.ID] = ci
	}
}

func (set *CaseSet) FilterByFunc(fn FuncCaseMatch) *CaseSet {
	out := NewCaseSet()
	if fn == nil {
		return out
	}
	for _, c := range set.Cases {
		if fn(c) {
			out.Cases[c.ID] = c
		}
	}
	return out
}

func (set *CaseSet) Get(caseID uint) (Case, bool) {
	c, ok := set.Cases[caseID]
	return c, ok
}

func (set *CaseSet) IDs() []uint { return maputil.Keys(set.Cases) }

func (set *CaseSet) IDsByFunc(fn FuncCaseMatch) []uint {
	var out []uint
	for _, c := range set.Cases {
		if fn(c) {
			out = append(out, c.ID)
		}
	}
	out = slicesutil.Dedupe(out)
	slices.Sort(out)
	return out
}

func (set *CaseSet) IDsBySection(sectionID uint) []uint {
	return set.IDsByFunc(
		func(c Case) bool {
			if c.SectionID != nil && pointer.Dereference(c.SectionID) == sectionID {
				return true
			} else {
				return false
			}
		},
	)
}

func (set *CaseSet) Len() uint { return uint(len(set.Cases)) }

func (set *CaseSet) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return jsonutil.WriteFile(filename, set, prefix, indent, perm)
}
