package gotestrail

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/sort/sortutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
)

type Section struct {
	ID           uint      `json:"id"`
	SuiteID      uint      `json:"suite_id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	ParentID     *uint     `json:"parent_id"`
	DisplayOrder uint      `json:"display_order"`
	Depth        uint      `json:"depth"`
	ChildIDs     []uint    `json:"child_ids"` // not in API
	Lineage      Metadatas `json:"lineage"`
}

type SectionSet struct {
	Sections map[uint]Section `json:"sections,omitempty"`
}

func NewSectionSet() *SectionSet {
	return &SectionSet{Sections: map[uint]Section{}}
}

func ReadFileSectionSet(filename string) (*SectionSet, error) {
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		ss := &SectionSet{}
		return ss, json.Unmarshal(b, ss)
	}
}

func (set *SectionSet) Add(sections ...Section) {
	for _, s := range sections {
		set.Sections[s.ID] = s
	}
}

func (set *SectionSet) GetChildren(id uint) []Section {
	var out []Section
	for _, s := range set.Sections {
		if s.ParentID != nil && pointer.Dereference(s.ParentID) == id {
			out = append(out, s)
		}
	}
	return out
}

// GetChildrenIDsFlat returns a list of sectionIDs in flat sorted order, e.g. not depth first order.
func (set *SectionSet) GetChildrenIDsFlat(id uint, recurse bool) ([]uint, error) {
	var ids []uint
	if s, ok := set.Sections[id]; !ok {
		return ids, fmt.Errorf("section not found (%d)", id)
	} else {
		if len(s.ChildIDs) > 0 {
			ids = append(ids, s.ChildIDs...)
			if recurse {
				for _, childID := range s.ChildIDs {
					if moreIDs, err := set.GetChildrenIDsFlat(childID, recurse); err != nil {
						return ids, err
					} else if len(moreIDs) > 0 {
						ids = append(ids, moreIDs...)
					}
				}
			}
			ids = slicesutil.Dedupe(ids)
			slices.Sort(ids)
		}
		return ids, nil
	}
}

func (set *SectionSet) GetByName(name string, depth int) []Section {
	var out []Section
	for _, s := range set.Sections {
		if s.Name != name {
			continue
		} else if depth >= 0 && s.Depth != uint(depth) {
			continue
		} else {
			out = append(out, s)
		}
	}
	return out
}

func (set *SectionSet) Inflate() error {
	for childID, childSection := range set.Sections {
		if childSection.ParentID == nil {
			continue
		}
		if err := set.addParentChildMapping(pointer.Dereference(childSection.ParentID), childID); err != nil {
			return err
		}
		childSection.Lineage = set.buildLineage(childID)
		set.Sections[childID] = childSection
	}
	return nil
}

func (set *SectionSet) buildLineage(childID uint) Metadatas {
	var lineage Metadatas
	curChildID := childID
	for {
		if parentID, ok := set.getParentID(curChildID); ok {
			parentName := ""
			if parentSection, ok := set.Sections[parentID]; ok {
				parentName = parentSection.Name
			}
			lineage = append(lineage, Metadata{
				ID:   parentID,
				Name: parentName})
			curChildID = parentID
		} else {
			break
		}
	}
	slices.Reverse(lineage)
	return lineage
}

func (set *SectionSet) getParentID(childID uint) (uint, bool) {
	if childSection, ok := set.Sections[childID]; !ok {
		return 0, false
	} else if childSection.ParentID == nil {
		return 0, false
	} else if parentSection, ok := set.Sections[pointer.Dereference(childSection.ParentID)]; !ok {
		return 0, false
	} else {
		return parentSection.ID, true
	}
}

func (set *SectionSet) addParentChildMapping(parentID, childID uint) error {
	parentSection, ok := set.Sections[parentID]
	if !ok {
		return errors.New("parent section not found")
	}
	parentSection.ChildIDs = append(parentSection.ChildIDs, childID)
	if len(parentSection.ChildIDs) > 1 {
		parentSection.ChildIDs = slicesutil.Dedupe(parentSection.ChildIDs)
		sortutil.Slice(parentSection.ChildIDs)
	}
	set.Sections[parentID] = parentSection
	return nil
}

func (set *SectionSet) IDs() []uint { return maputil.Keys(set.Sections) }
func (set *SectionSet) Len() uint   { return uint(len(set.Sections)) }

func (set *SectionSet) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return jsonutil.WriteFile(filename, set, prefix, indent, perm)
}
