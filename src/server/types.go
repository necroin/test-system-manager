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

type TestPlanDescriptor struct {
	Description *string                   `json:"description,omitempty"`
	TestCases   []TestPlanDescriptorCases `json:"cases,omitempty"`
}

type ProjectUserUpdateRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
