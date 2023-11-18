package server

import (
	"encoding/json"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/settings"
)

func (server *Server) ProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	server.PageHandler(responseWriter, settings.InterfaceProjectsHTML)
}

func (server *Server) ProjectsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	projectsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "Projects",
		Fields: []dbi.Field{
			{
				Name: "Id",
			},
			{
				Name: "Name",
			},
		},
	})

	if projectsResponse.Error != nil {
		projectsResponse.Success = false
		json.NewEncoder(responseWriter).Encode(projectsResponse)
		return
	}

	for _, record := range projectsResponse.Records {
		projectId := record.Fields["Id"]

		testCaseResponse := server.db.SelectRequest(&dbi.Request{
			Table: "TSM_TestCase",
			Fields: []dbi.Field{
				{
					Name: "Count(*)",
				},
			},
			Filters: []dbi.Filter{
				{
					Name:     "ProjectId",
					Operator: "=",
					Value:    projectId,
				},
			},
		})

		if testCaseResponse.Error != nil {
			projectsResponse.Success = false
			json.NewEncoder(responseWriter).Encode(projectsResponse)
			return
		}

		record.Fields["TestCaseCount"] = testCaseResponse.Records[0].Fields["Count(*)"]
	}

	json.NewEncoder(responseWriter).Encode(projectsResponse)
}
