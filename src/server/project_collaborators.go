package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"
	"tsm/src/settings"

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
	params := mux.Vars(request)
	projectId := params["id"]

	data, _ := io.ReadAll(request.Body)
	username := string(data)

	server.db.InsertRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Fields: []dbi.Field{
			{
				Name:  "Username",
				Value: fmt.Sprintf("'%s'", username),
			},
			{
				Name:  "ProjectId",
				Value: projectId,
			},
			{
				Name:  "Role",
				Value: fmt.Sprintf("'%s'", settings.ProjectRoleGuest),
			},
		},
	})
}

func (server *Server) ProjectCollaboratorsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	data, _ := io.ReadAll(request.Body)
	username := string(data)

	server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Filters: []dbi.Filter{
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", username),
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})
}

func (server *Server) ProjectCollaboratorsUpdateHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	updateRequest := &ProjectUserUpdateRequest{}
	json.NewDecoder(request.Body).Decode(updateRequest)

	server.db.UpdateRequest(&dbi.Request{
		Table: "TSM_ProjectUsers",
		Fields: []dbi.Field{{
			Name:  "Role",
			Value: fmt.Sprintf("'%s'", updateRequest.Role),
		}},
		Filters: []dbi.Filter{
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", updateRequest.Username),
			},
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    projectId,
			},
		},
	})
}
