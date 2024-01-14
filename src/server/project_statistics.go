package server

import (
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/settings"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectStatisticsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	token := params["token"]
	projectId := params["id"]

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
	// params := mux.Vars(request)
	// projectId := params["id"]
}
