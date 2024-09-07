package gotestrail

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/grokify/mogo/strconv/strconvutil"
)

func ReadFileCaseFields(filename string) (CaseFields, error) {
	if b, err := os.ReadFile(filename); err != nil {
		return CaseFields{}, err
	} else {
		return ReadCaseFields(b)
	}
}

func ReadCaseFields(b []byte) (CaseFields, error) {
	var out []CaseField
	return out, json.Unmarshal(b, &out)
}

type CaseFields []CaseField

type CaseField struct {
	ID         int               `json:"id"`
	IsActive   bool              `json:"is_active"`
	Name       string            `json:"name"`
	SystemName string            `json:"system_name"`
	Label      string            `json:"label"`
	Configs    []CaseFieldConfig `json:"configs"`
}

type CaseFieldConfig struct {
	Context CaseFieldContext       `json:"context"`
	ID      string                 `json:"id"`
	Options CaseFieldConfigOptions `json:"options"`
}

type CaseFieldContext struct {
	IsGlobal   bool  `json:"is_global"`
	ProjectIDs []int `json:"project_ids"`
}

type CaseFieldConfigOptions struct {
	IsRequired bool   `json:"is_required"`
	Items      string `json:"items"`
	Format     string `json:"format"`
}

var rxItem = regexp.MustCompile(`^(\d+),\s+(\S.+)$`)

func (opts CaseFieldConfigOptions) ItemsMap() (map[uint]string, error) {
	return parseCaseFieldItemsMap(opts.Items)
}

func parseCaseFieldItemsMap(s string) (map[uint]string, error) {
	out := map[uint]string{}
	s = strings.TrimSpace(s)
	if s == "" {
		return out, nil
	}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		m := rxItem.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}
		num := strings.TrimSpace(m[1])
		name := strings.TrimSpace(m[2])
		if num == "" || name == "" {
			return out, fmt.Errorf("cannot parse items map (%s)", s)
		}
		numUint, err := strconvutil.Atou(num)
		if err != nil {
			return out, fmt.Errorf("items map key is not a number (%s)", num)
		} else {
			out[numUint] = name
		}
	}
	return out, nil
}

type CaseFieldSet struct {
	Data map[int]CaseField
}
