package server

import (
	"encoding/json"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectStatisticsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

	if server.GetUserProjectRole(token, projectId) < roleCreator {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	projectsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "Id",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if projectsResponse.Error != nil {
		logger.Error(projectsResponse.Error.Error())
		responseWriter.Write([]byte(projectsResponse.Error.Error()))
		return
	}

	PageDescriptor := PageDescriptor{
		ProjectId:   projectId,
		ProjectName: projectsResponse.Records[0].Fields["Name"],
	}
	server.PageHandler(responseWriter, settings.InterfaceProjectStatisticsHTML, PageDescriptor, token)
}

func (server *Server) ProjectStatisticsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

	if server.GetUserProjectRole(token, projectId) < roleCreator {
		responseWriter.Write([]byte("Permission denied"))
		return
	}

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Stat",
		Fields: []dbi.Field{
			{Name: "TestPlanId"},
			{Name: "TestCaseId"},
			{Name: "TestRunId"},
			{Name: "Result"},
			{Name: "Datetime"},
		},
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})

	if response.Error != nil {
		logger.Error(response.Error.Error())
		responseWriter.Write([]byte(response.Error.Error()))
		return
	}

	json.NewEncoder(responseWriter).Encode(response)
}
