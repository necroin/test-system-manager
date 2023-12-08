package server

import (
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectPlanHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]
	testPlanId := params["planId"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_TestPlan",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    testPlanId,
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
		TestPlanId:   testPlanId,
		TestPlanName: response.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfacePlanHTML, PageDescriptor)
}

func (server *Server) ProjectPlanSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
}
