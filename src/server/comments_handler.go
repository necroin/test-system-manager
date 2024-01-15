package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tsm/src/db/dbi"
	"tsm/src/logger"

	"github.com/gorilla/mux"
)

func (server *Server) ProjectCommentsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["projid"]
	objectId := params["objId"]
	objectType := params["type"]

	projectCommentsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Comments",
		Fields: []dbi.Field{
			{
				Name: "Username",
			},
			{
				Name: "Content",
			},
			{
				Name: "Id",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", objectId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", objectType),
			},
		},
	})

	if projectCommentsResponse.Error != nil {
		logger.Error("%s", projectCommentsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectCommentsResponse)
		return
	}

	json.NewEncoder(responseWriter).Encode(projectCommentsResponse)
}

func (server *Server) ProjectCommentsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["projid"]
	objectId := params["objId"]
	objectType := params["type"]
	token := params["token"]
	content, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	logger.Debug("123", projectId, objectId, objectType)
	response := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_Comments",
		Fields: []dbi.Field{
			{
				Name:  "Id",
				Value: fmt.Sprintf("(select COALESCE((MAX(Id) + 1), 1) from TSM_Comments where ProjectId = %s and ObjectId = %s and ObjectType = '%s')", projectId, objectId, objectType),
			},
			{
				Name:  "Content",
				Value: fmt.Sprintf("'%s'", string(content)),
			},
			{
				Name:  "ProjectId",
				Value: fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:  "ObjectId",
				Value: fmt.Sprintf("'%s'", objectId),
			},
			{
				Name:  "ObjectType",
				Value: fmt.Sprintf("'%s'", objectType),
			},
			{
				Name:  "Username",
				Value: fmt.Sprintf("'%s'", server.FindUsernameByToken(token)),
			},
		},
		Filters: []dbi.Filter{},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) ProjectCommentsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id := params["id"]
	projectId := params["projid"]
	objectId := params["objId"]
	objectType := params["type"]
	token := params["token"]
	response := server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Comments",
		Filters: []dbi.Filter{
			{
				Name:     "ProjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", objectId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", objectType),
			},
			{
				Name:     "Id",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", id),
			},
			{
				Name:     "Username",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", server.FindUsernameByToken(token)),
			},
		},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}
