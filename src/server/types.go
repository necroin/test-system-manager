package server

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
}

type SearchRequest struct {
	SearchFilter string `json:"search"`
}

type TestCaseDescriptor struct {
	Description *string `json:"description,omitempty"`
	Scenario    *string `json:"scenario,omitempty"`
}

type TestPlanDescriptor struct {
	Description *string `json:"description,omitempty"`
	TestCases   *string `json:"scenario,omitempty"`
}
