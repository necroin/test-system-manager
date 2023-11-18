package server

type TestCaseDescriptor struct {
	Description *string `json:"description,omitempty"`
	Scenario    *string `json:"scenario,omitempty"`
}
