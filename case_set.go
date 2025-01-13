package gotestrail

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"slices"

	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/os/osutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
)

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

func (set *CaseSet) ReadFileJSONs(filenames ...string) error {
	for _, fn := range filenames {
		if err := set.ReadFileJSON(fn); err != nil {
			return err
		}
	}
	return nil
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

func (set *CaseSet) ReadDirJSONRawAPIs(dir string, rx *regexp.Regexp) error {
	entries, err := osutil.ReadDirMore(dir, rx, false, true, false)
	if err != nil {
		return err
	}
	filenames := entries.Names(dir)
	return set.ReadFileJSONRawAPIs(filenames...)
}

func (set *CaseSet) ReadFileJSONRawAPIs(filenames ...string) error {
	for _, fn := range filenames {
		if err := set.ReadFileJSONRawAPI(fn); err != nil {
			return err
		}
	}
	return nil
}

func (set *CaseSet) ReadFileJSONRawAPI(filename string) error {
	if new, err := ReadFileAPIResponseGetCases(filename); err != nil {
		return err
	} else {
		for _, ci := range new.Cases {
			set.Cases[ci.ID] = ci
		}
	}
	return nil
}

func (set *CaseSet) Add(c ...Case) {
	for _, ci := range c {
		set.Cases[ci.ID] = ci
	}
}

var ErrResponseCannotBeNil = errors.New("response cannot be nil")

/*
func (set *CaseSet) ProcCaseResponse(r *http.Response) (string, error) {
	if r == nil {
		return "", ErrResponseCannotBeNil
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

}
*/

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

func (set *CaseSet) LineageStringsHistogram(sectionSet *SectionSet, strSep string) (*histogram.Histogram, error) {
	if sectionSet == nil {
		return nil, errors.New("gotestrail.SectionSet must be supplied")
	}
	h := histogram.NewHistogram("")
	for _, c := range set.Cases {
		secID := c.SectionID
		if secID == nil {
			return nil, errors.New("section id cannot be nil")
		}
		lin := sectionSet.BuildLineage(pointer.Dereference(secID))
		strs := lin.LineageIDsStrings()
		for _, str := range strs {
			strNames, err := sectionSet.LineageIDsStringToNames(str, " > ")
			if err != nil {
				return nil, err
			}
			h.Add(strNames, 1)
		}
	}
	return h, nil
}

func (set *CaseSet) Len() uint { return uint(len(set.Cases)) }

func (set *CaseSet) WriteFileJSON(filename string, perm os.FileMode, prefix, indent string) error {
	return jsonutil.MarshalFile(filename, set, prefix, indent, perm)
}
