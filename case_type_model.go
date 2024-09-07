package gotestrail

import (
	"encoding/json"
	"os"
)

type CaseType struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

type CaseTypeSet struct {
	CaseTypes map[uint]CaseType `json:"caseTypes"`
}

func NewCaseTypeSet() *CaseTypeSet {
	return &CaseTypeSet{CaseTypes: map[uint]CaseType{}}
}

func (set *CaseTypeSet) Add(items ...CaseType) {
	if set.CaseTypes == nil {
		set.CaseTypes = map[uint]CaseType{}
	}
	for _, item := range items {
		set.CaseTypes[item.ID] = item
	}
}

func ReadFileCaseTypeSet(filename string) (*CaseTypeSet, error) {
	set := NewCaseTypeSet()
	if b, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return set, json.Unmarshal(b, set)
	}
}

func (set *CaseTypeSet) ReadFileJSON(filename string) error {
	if new, err := ReadFileCaseTypeSet(filename); err != nil {
		return err
	} else {
		for k, v := range new.CaseTypes {
			set.CaseTypes[k] = v
		}
	}
	return nil
}
