package gotestrail

import (
	"fmt"
	"strings"
)

type Metadatas []Metadata

type Metadata struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (mds Metadatas) LineageIDsString() string {
	var parts []string
	for _, md := range mds {
		parts = append(parts, fmt.Sprintf("%d", md.ID))
	}
	return strings.Join(parts, ".")
}

func (mds Metadatas) LineageIDsStrings() []string {
	var strs []string
	var parts []string
	for _, md := range mds {
		parts = append(parts, fmt.Sprintf("%d", md.ID))
		strs = append(strs, strings.Join(parts, "."))
	}
	return strs
}
