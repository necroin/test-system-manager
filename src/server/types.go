package server

type SearchRequest struct {
	SearchFilter string `json:"search"`
}

type TestCaseDescriptor struct {
	Description *string `json:"description,omitempty"`
	Scenario    *string `json:"scenario,omitempty"`
}
