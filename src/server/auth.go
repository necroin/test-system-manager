package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"tsm/src/db/dbi"
	"tsm/src/settings"
)

func (server *Server) AuthHandler(responseWriter http.ResponseWriter, request *http.Request) {
	PageDescriptor := PageDescriptor{
		Style:  settings.GetAuthenticationStyle(),
		Script: settings.GetAuthenticationScript(),
	}
	pageTemplate, _ := template.New("HtmlPage").Parse(settings.GetAuthenticationTemlate())
	pageTemplate.Execute(responseWriter, PageDescriptor)
}

func (server *Server) AuthRegisterHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &ClientAuth{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Users",
		Fields: []dbi.Field{
			{
				Name: "Token",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", data.Username),
			},
			{
				Name:     "Password",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", data.Password),
			},
		},
	})

	if response.Error != nil {
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	if len(response.Records) > 0 {
		return
	}

	token := hex.EncodeToString([]byte(data.Username + ";" + data.Password))

	server.db.InsertRequest(&dbi.Request{
		Table: "TSM_Users",
		Fields: []dbi.Field{
			{
				Name:  "Username",
				Value: fmt.Sprintf("'%s'", data.Username),
			},
			{
				Name:  "Password",
				Value: fmt.Sprintf("'%s'", data.Password),
			},
			{
				Name:  "Token",
				Value: fmt.Sprintf("'%s'", token),
			},
		},
	})
	responseWriter.Write([]byte(token))
}

func (server *Server) AuthTokenHandler(responseWriter http.ResponseWriter, request *http.Request) {
	data := &ClientAuth{}
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Users",
		Fields: []dbi.Field{
			{
				Name: "Token",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", data.Username),
			},
			{
				Name:     "Password",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", data.Password),
			},
		},
	})

	if response.Error != nil {
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	if len(response.Records) == 0 {
		return
	}

	token := response.Records[0].Fields["Token"]
	responseWriter.Write([]byte(token))
}
