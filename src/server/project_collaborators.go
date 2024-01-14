package server

import (
	"encoding/json"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectCollaboratorsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	response := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Fields: []dbi.Field{
			{
				Name: "Username",
			},
			{
				Name: "Role",
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

	if response.Error != nil {
		logger.Error("%s", response.Error)
		json.NewEncoder(responseWriter).Encode(response)
		return
	}

	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) ProjectCollaboratorsAddHandler(responseWriter http.ResponseWriter, request *http.Request) {

}

func (server *Server) ProjectCollaboratorsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {

}
