package gotestrail

import (
	"strings"

	"github.com/grokify/mogo/pointer"
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
	XLineageSectionIDs   string  `json:"x_lineage_section_ids"`
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
