package gotestrail

const (
	IndexPath                   = "index.php?"
	APIPathCasesGet             = "/api/v2/get_cases/"
	APIPathCasesGetProjectID    = "/api/v2/get_cases/%d&limit=%d&offset=%d"
	APIPathCaseTypesGet         = "/api/v2/get_case_types/"
	APIPathSectionsGet          = "/api/v2/get_sections/"
	APIPathSectionsGetProjectID = "/api/v2/get_sections/%d&limit=%d&offset=%d"

	LimitMax uint = 250

	ParamSuiteID = "suite_id"

	ParamLimit  = "limit"
	ParamOffset = "offset"
)
