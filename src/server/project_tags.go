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

func (server *Server) ProjectTagsSelectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	projectTagsResponse := server.db.SelectRequest(&dbi.Request{
		Table: "TSM_Tags",
		Fields: []dbi.Field{
			{
				Name: "Name",
			},
		},
		Filters: []dbi.Filter{
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    "'Project'",
			},
		},
	})

	if projectTagsResponse.Error != nil {
		logger.Error("%s", projectTagsResponse.Error)
		json.NewEncoder(responseWriter).Encode(projectTagsResponse)
		return
	}

	json.NewEncoder(responseWriter).Encode(projectTagsResponse)
}

func (server *Server) ProjectTagsInsertHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	tagName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.InsertRequest(&dbi.Request{
		Table: "TSM_Tags",
		Fields: []dbi.Field{
			{
				Name:  "Name",
				Value: fmt.Sprintf("'%s'", string(tagName)),
			},
			{
				Name:  "ObjectId",
				Value: fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:  "ObjectType",
				Value: "'Project'",
			},
		},
		Filters: []dbi.Filter{},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}

func (server *Server) ProjectTagsDeleteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	projectId := params["id"]

	tagName, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("%s", err)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	response := server.db.DeleteRequest(&dbi.Request{
		Table: "TSM_Tags",
		Filters: []dbi.Filter{
			{
				Name:     "Name",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", string(tagName)),
			},
			{
				Name:     "ObjectId",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", projectId),
			},
			{
				Name:     "ObjectType",
				Operator: "=",
				Value:    fmt.Sprintf("'%s'", settings.ObjectTypeProject),
			},
		},
	})

	if response.Error != nil {
		logger.Error("%s", response.Error)
	}
	json.NewEncoder(responseWriter).Encode(response)
}
