package dbi

type Filter struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Record struct {
	Fields map[string]string `json:"fields"`
}

type Request struct {
	Table   string   `json:"table"`
	Fields  []Field  `json:"fields"`
	Filters []Filter `json:"filters"`
}

type Response struct {
	Records []Record `json:"records,omitempty"`
	Success bool     `json:"success"`
	Error   error    `json:"error,omitempty"`
}

type Database interface {
	InsertRequest(request *Request) *Response
	SelectRequest(request *Request) *Response
	DeleteRequest(request *Request) *Response
	UpdateRequest(request *Request) *Response
}
