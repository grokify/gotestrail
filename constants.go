package gotestrail

const (
	IndexPath                   = "index.php?"
	APIPathCasesGet             = "/api/v2/get_cases/"
	APIPathCasesGetProjectID    = "/api/v2/get_cases/%d&limit=%d&offset=%d"
	APIPathCaseFieldsGet        = "/api/v2/get_case_fields/"
	APIPathCaseTypesGet         = "/api/v2/get_case_types/"
	APIPathSectionsGet          = "/api/v2/get_sections/"
	APIPathSectionsGetProjectID = "/api/v2/get_sections/%d&limit=%d&offset=%d"

	LimitMax uint = 250

	QueryParamLimit     = "limit"
	QueryParamOffset    = "offset"
	QueryParamSectionID = "section_id"
	QueryParamSuiteID   = "suite_id"

	SlugCase     = "case"
	SlugCaseType = "case_type"
	SlugSection  = "section"
)
