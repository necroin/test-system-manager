package server

type ClientAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PageDescriptor struct {
	Url          string
	Style        string
	Script       string
	ProjectId    string
	ProjectName  string
	TestCaseId   string
	TestCaseName string
	TestPlanId   string
	TestPlanName string
	TestRunId    string
}

type SearchRequest struct {
	SearchFilter string `json:"search"`
}

type TestCaseDescriptor struct {
	Description *string `json:"description,omitempty"`
	Scenario    *string `json:"scenario,omitempty"`
}

type TestPlanDescriptorCases struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TestRunDescriptorCases struct {
	TestRunId  string `json:"run_id"`
	TestCaseId string `json:"case_id"`
	Result     string `json:"result"`
	Datetime   string `json:"datetime"`
	Comment    string `json:"comment"`
}

type TestPlanDescriptor struct {
	Description *string                   `json:"description,omitempty"`
	TestCases   []TestPlanDescriptorCases `json:"cases,omitempty"`
	TestRuns    []TestRunDescriptorCases   `json:"runs,omitempty"`
}

type ProjectUserUpdateRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
type TestPlanRunDescriptorCases struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Scenario    string `json:"scenario"`
	Description string `json:"case_description"`
}

type TestPlanRunDescriptor struct {
	//Id string `json:"run_id"`
	Description *string                      `json:"description,omitempty"`
	TestCases   []TestPlanRunDescriptorCases `json:"cases,omitempty"`
}

type TestCaseStatDescriptor struct {
	Result   string `json:"result"`
	Id       string `json:"case_id"`
	RunId    string `json:"run_id"`
	DateTime string `json:"datetime"`
	Comment  string `json:"comment"`
}

type TestPlanStatDescriptor struct {
	TestCases []TestCaseStatDescriptor `json:"cases,omitempty"`
}



// type StatCaseDescriptor struct {
// 	TestCaseId string `json:"case_id"`
// 	Result     string `json:"result"`
// }

// type StatItemDescriptor struct {
// 	TestPlanId string               `json:"plan_id"`
// 	TestRunId  string               `json:"run_id"`
// 	Datetime   string               `json:"datetime"`
// 	Cases      []StatCaseDescriptor `json:"stat_test_cases"`
// }

// type StatDescriptor struct {
// 	Items []StatItemDescriptor `json:"items,omitempty"`
// }
