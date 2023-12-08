package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectCaseHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testCaseId := params["caseId"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestCase",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    testCaseId,
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	PageDescriptor := PageDescriptor{
		ProjectId:    projectId,
		TestCaseId:   testCaseId,
		TestCaseName: response.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfaceCaseHTML, PageDescriptor)
}

func (server *Server) ProjectCaseSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testCaseId := params["caseId"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestCase",
		Fields: []dbi.Field{
			{
				Name: "Description",
			},
			{
				Name: "Scenario",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    testCaseId,
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error)
	}

	if len(response.Records) == 0 {
		return
	}

	description := response.Records[0].Fields["Description"]
	scenario := response.Records[0].Fields["Scenario"]

	data := &TestCaseDescriptor{
		Description: &description,
		Scenario:    &scenario,
	}
	json.NewEncoder(responseWriter).Encode(data)
}

func (server *Server) ProjectCaseUpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testCaseId := params["caseId"]
	fields := []dbi.Field{}

	data := &TestCaseDescriptor{}
	json.NewDecoder(request.Body).Decode(data)

	if data.Description != nil {
		fields = append(fields, dbi.Field{
			Name:  "Description",
			Value: fmt.Sprintf("'%s'", *data.Description),
		})
	}

	if data.Scenario != nil {
		fields = append(fields, dbi.Field{
			Name:  "Scenario",
			Value: fmt.Sprintf("'%s'", *data.Scenario),
		})
	}

	if len(fields) > 0 {
		response := server.db.UpdateRequest(&dbi.Request{
			Table:  "TSM_TestCase",
			Fields: fields,
			Filters: []dbi.Filter{
				{
					Name:     "Id",
					Operator: "=",
					Value:    testCaseId,
				},
				{
					Name:     "ProjectId",
					Operator: "=",
					Value:    projectId,
				},
			},
		})

		if response.Error != nil {
			logger.Error(response.Error)
		}
	}
}
